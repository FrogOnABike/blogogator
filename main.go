package main

import (
	"fmt"

	"github.com/frogonabike/blogogator/internal/config"
)

func main() {
	configFile, err := config.Read()
	if err != nil {
		return
	}
	fmt.Println(configFile)

	config.SetUser("mark")

	configFileUpdated, err := config.Read()
	if err != nil {
		return
	}
	fmt.Println(configFileUpdated)
}
