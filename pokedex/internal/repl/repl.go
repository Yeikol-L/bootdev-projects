package repl

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/yeikol-l/bootdev/pokedex/internal/api"
	"github.com/yeikol-l/bootdev/pokedex/internal/pokecache"
)

type Pokedex struct {
	pokemons []api.PokemonDetails
}

func (p Pokedex) find(name string) (pokemon api.PokemonDetails, ok bool) {
	for _, v := range p.pokemons {
		if v.Name == name {
			return v, true
		}
	}
	return api.PokemonDetails{}, false
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}
type Program struct {
	scanner *bufio.Scanner
	cfg     config
}
type config struct {
	commands  map[string]cliCommand
	mapCursor int
	cache     *pokecache.Cache
	pokedex   Pokedex
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "it displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "it displays the names of the previous 20 location areas in the Pokemon world",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "it displays the pokemons in the location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "it tries to catch the pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "it displays a catched pokemon details",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "it displays all catched pokemon names",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func NewProgram(r io.Reader) Program {
	return Program{
		scanner: bufio.NewScanner(r),
		cfg:     config{commands: getCommands(), mapCursor: -20, cache: pokecache.NewCache(time.Second * 15)},
	}
}

func (p Program) Start() {
	for true {
		fmt.Print("Pokedex > ")
		p.scanner.Scan()
		input := cleanInput(p.scanner.Text())
		command, ok := p.cfg.commands[input[0]]
		if !ok {
			fmt.Print("Unknown command\n")
		} else {
			err := command.callback(&p.cfg, input[1:])
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}

func commandExit(c *config, params []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, params []string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range c.commands {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(c *config, params []string) error {
	c.mapCursor += 20
	locationAreas, err := api.GetLocationAreas(c.cache, c.mapCursor)
	if err != nil {
		return err
	}
	for _, a := range locationAreas {
		fmt.Println(a.Name)
	}
	return nil
}
func commandMapb(c *config, params []string) error {
	if c.mapCursor > 0 {
		c.mapCursor -= 20
	}
	locationAreas, err := api.GetLocationAreas(c.cache, c.mapCursor)
	if err != nil {
		return err
	}
	for _, a := range locationAreas {
		fmt.Println(a.Name)
	}
	return nil
}
func commandExplore(c *config, params []string) error {
	fmt.Printf("Exploring %s...\n", params[0])
	pokemons, err := api.GetLocationAreaPokemons(c.cache, params[0])
	if err != nil {
		return err
	}
	fmt.Printf("Found pokemons: \n")
	for _, p := range pokemons {
		fmt.Println("- " + p.Name)
	}
	return nil
}

func commandCatch(c *config, params []string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", params[0])
	pokemon, err := api.GetPokemonDetail(c.cache, params[0])
	if err != nil {
		return err
	}
	difficulty := pokemon.BaseExperience / 50
	if pokemon.BaseExperience > rand.Intn(pokemon.BaseExperience*difficulty) {
		c.pokedex.pokemons = append(c.pokedex.pokemons, pokemon)
		fmt.Printf("%s was caught!\n", params[0])
	} else {
		fmt.Printf("Failed to catch %s\n", params[0])
	}

	return nil
}
func commandInspect(c *config, params []string) error {
	pokemon, ok := c.pokedex.find(params[0])
	if !ok {
		fmt.Printf("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("- Name: %s\n", pokemon.Name)
	fmt.Printf("- Height: %d\n", pokemon.Height)
	fmt.Printf("- Weight: %d\n", pokemon.Weight)
	return nil
}
func commandPokedex(c *config, _ []string) error {
	if len(c.pokedex.pokemons) == 0 {
		fmt.Println("Your pokedex is empty:")
		return nil
	}
	fmt.Println("Your pokedex:")
	for _, v := range c.pokedex.pokemons {
		fmt.Printf("- %s\n", v.Name)
	}

	return nil
}

func cleanInput(text string) []string {
	return strings.Split(strings.ToLower(strings.Trim(text, " ")), " ")
}
