package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/api"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/config"
	"github.com/yeikol-l/bootdev-projects/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	wrkdir, err := os.Getwd()
	if err != nil {
		fmt.Println("Could get working dir")
		return
	}
	godotenv.Load(wrkdir + "/.env")
	serverConfig, err := config.NewApiConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Could not connect to db")
		return
	}
	dbQueries := database.New(db)
	api := api.New(dbQueries, serverConfig.IsDev(), serverConfig.GetJWTSecret(), serverConfig.GetPolkaKey())

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", serverConfig.GetPort()),
		Handler: api.Routes(),
	}
	fmt.Println("Running http server on port " + serverConfig.GetPort())
	defer db.Close()
	server.ListenAndServe()
}
