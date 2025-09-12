package main

import (
	"log"
	"os"

	"github.com/frogonabike/blogogator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	// Read the saved config file and save point the state struct to it
	configFile, err := config.Read()
	if err != nil {
		log.Fatalf("Unable to read config file:%v\n", err)
	}
	curState := state{
		config: &configFile,
	}
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
