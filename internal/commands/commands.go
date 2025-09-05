package commands

import (
	"errors"
	"fmt"

	"github.com/frogonabike/blogogator/internal/config"
)

type state struct {
	Config *config.Config
}

type command struct {
	Name string
	Args []string
}

func handlerLogin(s *state, cmd command) error {
	// Check that we only have a single username in the args slice, otherwise return an error
	if len(cmd.Args) != 1 {
		return errors.New("please enter a single username")
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Println("Username set successfully")
	return nil
}
