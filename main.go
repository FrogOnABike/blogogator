package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/frogonabike/blogogator/internal/config"
	"github.com/frogonabike/blogogator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db     *database.Queries
	config *config.Config
}

func main() {
	// Read the saved config file and save a pointer to the state struct
	configFile, err := config.Read()
	if err != nil {
		log.Fatalf("Unable to read config file:%v\n", err)
	}
	curState := state{
		config: &configFile,
	}
	// Pretty print the configFile for debugging
	fmt.Println(configFile)

	// Open connection to database
	db, err := sql.Open("postgres", curState.config.DbURL)
	if err != nil {
		log.Fatalf("Error connecting to database:%v\n", err)
	}
	// Create a new dbQueries and store in state
	dbQueries := database.New(db)
	curState.db = dbQueries
	fmt.Println(curState)

	// Initialise the command hanlders struct
	comHandlers := commands{Handlers: make(map[string]func(*state, command) error)}
	// Register commands
	comHandlers.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Too few arguments\n")
	}

	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	err = comHandlers.run(&curState, cmd)
	if err != nil {
		log.Fatalf("Unable to run command:%v", err)
	}

	// fmt.Println(configFile)

	// configFile.SetUser("mark")

	// updatedCfg, err := config.Read()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(updatedCfg)
}
