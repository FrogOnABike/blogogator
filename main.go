package main

import (
	"fmt"

	"github.com/frogonabike/blogogator/internal/config"
)

func main() {
	fmt.Println(config.Read())
}
