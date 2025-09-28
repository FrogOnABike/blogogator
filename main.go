package main

import (
	"database/sql"
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
	// fmt.Println(configFile)

	// Open connection to database
	db, err := sql.Open("postgres", curState.config.DbURL)
	if err != nil {
		log.Fatalf("Error connecting to database:%v\n", err)
	}

	// Create a new dbQueries and store in state
	dbQueries := database.New(db)
	curState.db = dbQueries
	// fmt.Println(curState)

	// Initialise the command hanlders struct
	comHandlers := commands{Handlers: make(map[string]func(*state, command) error)}

	// Register commands
	comHandlers.register("login", handlerLogin)
	comHandlers.register("register", handlerRegister)
	comHandlers.register("reset", handlerReset)
	comHandlers.register("users", handlerGetUsers)
	comHandlers.register("agg", handlerAgg)
	comHandlers.register("addfeed", handlerAddFeed)
	comHandlers.register("feeds", handlerFeeds)
	// ***Start of processing of user input***

	// Check we have at least a command passed
	args := os.Args
	if len(args) < 2 {
		log.Fatalf("Too few arguments\n")
	}
	// Parse the input: 0 - Always "gator", 1 - Command name, 2 > Arguments
	cmd := command{
		Name: args[1],
		Args: args[2:],
	}
	// Attempt to run the command, returning any errors if it's unable too be run
	err = comHandlers.run(&curState, cmd)
	if err != nil {
		log.Fatalf("Unable to run command:%v\n", err)
	}

	// fmt.Println(configFile)

	// configFile.SetUser("mark")

	// updatedCfg, err := config.Read()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(updatedCfg)
}
