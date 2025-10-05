package main

import (
	"errors"
	"fmt"
)

// Core command handling logic
type command struct {
	Name string
	Args []string
}

type commands struct {
	Handlers map[string]func(*state, command) error
}

// Method to run a given command, if it exists
func (c *commands) run(s *state, cmd command) error {
	handler, found := c.Handlers[cmd.Name]
	if !found {
		return errors.New("command not found")
	}
	err := handler(s, cmd)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	return nil
}

// Method to register a new command
func (c *commands) register(name string, f func(*state, command) error) {
	c.Handlers[name] = f
}
