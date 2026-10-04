package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/yeikol-l/bootdev/blog-aggregator/internal/config"
	"github.com/yeikol-l/bootdev/blog-aggregator/internal/database"
	"github.com/yeikol-l/bootdev/blog-aggregator/internal/rss"

	_ "github.com/lib/pq"
)

type state struct {
	config *config.Config
	db     *database.Queries
}
type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	handler, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("Handler %s not found", cmd.name)
	}
	return handler(s, cmd)
}
func (c *commands) register(name string, f func(*state, command) error) error {
	_, exists := c.handlers[name]
	if exists {
		return fmt.Errorf("Command already exists")
	}
	c.handlers[name] = f
	return nil
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
		return
	}
	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		fmt.Println(err)
		return
	}
	dbQueries := database.New(db)

	s := state{config: &cfg, db: dbQueries}
	cmds := commands{handlers: map[string]func(*state, command) error{}}
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAggregator)
	cmds.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cmds.register("feeds", handlerFeeds)
	cmds.register("follow", middlewareLoggedIn(handlerFollow))
	cmds.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	cmds.register("following", middlewareLoggedIn(handlerFollowing))
	cmds.register("browse", middlewareLoggedIn(handlerBrowse))

	args := os.Args
	if len(args) < 2 {
		fmt.Println("At least 1 argument required")
		os.Exit(1)
		return
	}
	err = cmds.run(&s, command{name: args[1], arguments: args[2:]})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user *database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		ctx := context.Background()
		user, err := s.db.GetUserByName(ctx, s.config.Username)
		if err != nil {
			return err
		}
		return handler(s, cmd, &user)
	}

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Expected one argument (username)")
	}
	ctx := context.Background()
	user, err := s.db.GetUserByName(ctx, cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("The user does not exist")
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Expected one argument (name)")
	}
	ctx := context.Background()
	_, err := s.db.GetUserByName(ctx, cmd.arguments[0])
	if err == nil {
		return fmt.Errorf("El usuario ya existe")
	}

	currentTime := time.Now()
	args := database.CreateUserParams{
		ID:        uuid.New(),
		Name:      cmd.arguments[0],
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	user, err := s.db.CreateUser(ctx, args)
	if err != nil {
		return err
	}
	err = s.config.SetUser(user.Name)
	if err != nil {
		return err
	}
	fmt.Printf("User has been created with id: %s\n", user.ID)
	return nil
}

func handlerReset(s *state, _ command) error {
	ctx := context.Background()
	err := s.db.DeleteAllUsers(ctx)
	if err != nil {
		return fmt.Errorf("No se puedierion eliminar los usuarios: %v", err)
	}
	fmt.Printf("Usuarios eliminados correctamente\n")
	return nil
}

func handlerUsers(s *state, _ command) error {
	ctx := context.Background()
	users, err := s.db.GetAllUsers(ctx)
	if err != nil {
		return err
	}
	for _, v := range users {
		fmt.Printf("* %s", v.Name)
		if v.Name == s.config.Username {
			fmt.Printf(" (current)")
		}
		fmt.Print("\n")
	}
	return nil
}
func handlerAggregator(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("At least one argument is required (time between requests) ej: 1s, 1m, 1h")
	}
	duration, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return err
	}
	return scrapeFeeds(s, duration)
}

func handlerAddFeed(s *state, c command, user *database.User) error {
	if len(c.arguments) < 2 {
		return fmt.Errorf("Se necesitan almenos dos argumentos (nombre, url)")
	}
	ctx := context.Background()
	currentTime := time.Now()
	args := database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      c.arguments[0],
		Url:       c.arguments[1],
		UserID:    user.ID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	feed, err := s.db.CreateFeed(ctx, args)
	if err != nil {
		return err
	}
	feedFollowArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	_, err = s.db.CreateFeedFollow(ctx, feedFollowArgs)
	if err != nil {
		return err
	}
	fmt.Printf("%v", feed)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	ctx := context.Background()
	feeds, err := s.db.GetAllFeeds(ctx)
	if err != nil {
		return err
	}
	for _, v := range feeds {
		fmt.Println("---------")
		fmt.Printf("id: %s\n", v.ID)
		fmt.Printf("name: %s\n", v.Name)
		fmt.Printf("url: %s\n", v.Url)
		fmt.Printf("user: %s\n", v.User)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user *database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("At leastt one argument required (url)")
	}
	ctx := context.Background()
	feed, err := s.db.GetFeedByUrl(ctx, cmd.arguments[0])
	if err != nil {
		return err
	}
	currentTime := time.Now()
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	follow, err := s.db.CreateFeedFollow(ctx, params)
	if err != nil {
		return err
	}
	fmt.Printf("Created follow: %s\n", follow.ID)
	fmt.Printf("Feed: %s\n", follow.FeedName)
	fmt.Printf("User: %s\n", follow.UserName)
	return nil
}

func handlerUnfollow(s *state, cmd command, user *database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("At leastt one argument required (url)")
	}
	ctx := context.Background()
	feed, err := s.db.GetFeedByUrl(ctx, cmd.arguments[0])
	if err != nil {
		return err
	}
	err = s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID})
	if err != nil {
		return err
	}
	fmt.Printf("Unfolowed feed: %s", feed.Name)
	return nil
}
func handlerFollowing(s *state, _ command, user *database.User) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	fmt.Println("Following:")
	for _, f := range feeds {
		fmt.Printf("- %s", f.FeedName)
	}
	return nil
}
func handlerBrowse(s *state, cmd command, user *database.User) error {
	limit := 5
	if len(cmd.arguments) > 0 {
		int, err := strconv.Atoi(cmd.arguments[0])
		if err == nil {
			limit = int
		}
	}
	ctx := context.Background()
	posts, err := s.db.GetUserPosts(ctx, database.GetUserPostsParams{UserID: user.ID, Limit: int32(limit)})
	if err != nil {
		return err
	}
	fmt.Println("User posts:")
	for _, f := range posts {
		fmt.Printf("Title: %s\n", f.Title)
		fmt.Printf("Description: %s\n", f.Description)
		fmt.Printf("Published at: %s", f.PublishedAt)
	}
	return nil
}

func scrapeFeeds(s *state, time_between_req time.Duration) error {
	ctx := context.Background()
	ticker := time.NewTicker(time_between_req)
	for ; ; <-ticker.C {
		nextFeed, err := s.db.GetNextFeedToFetch(ctx)
		if err != nil {
			return err
		}
		err = s.db.MarkFeedFetched(ctx, nextFeed.ID)
		if err != nil {
			return err
		}
		feed, err := rss.FetchFeed(ctx, nextFeed.Url)
		if err != nil {
			return err
		}
		fmt.Printf("Iterando items de: %s", feed.Channel.Title)
		for _, i := range feed.Channel.Item {
			parsedDate, err := time.Parse(time.RFC1123Z, i.PubDate)
			currentTime := time.Now()
			if err != nil {
				parsedDate = currentTime
			}
			params := database.CreatePostParams{
				ID:          uuid.New(),
				Title:       i.Title,
				Url:         i.Link,
				FeedID:      nextFeed.ID,
				Description: i.Description,
				PublishedAt: parsedDate,
				CreatedAt:   currentTime,
				UpdatedAt:   currentTime,
			}
			s.db.CreatePost(ctx, params)
		}
	}
}
