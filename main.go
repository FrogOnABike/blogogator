package main

import (
	"fmt"

	"github.com/frogonabike/blogogator/internal/config"
)

func main() {
	configFile, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(configFile)

	config.SetUser("mark")

	updatedCfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(updatedCfg)
}
