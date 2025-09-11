package main

import (
	"fmt"

	"github.com/frogonabike/blogogator/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	configFile, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	curState := state{
		config: &configFile,
	}

	comHandlers := commands(
		"login", handlerLogin,
	)
	// fmt.Println(configFile)

	// configFile.SetUser("mark")

	// updatedCfg, err := config.Read()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(updatedCfg)
}
