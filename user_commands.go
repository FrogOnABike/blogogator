package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"

	"github.com/frogonabike/blogogator/internal/database"
)

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

// ***Define handler functions of commands below***

// Login command
func handlerLogin(s *state, cmd command) error {
	// Check that we only have a single username in the args slice, otherwise return an error
	if len(cmd.Args) != 1 {
		log.Fatalf("please enter a username")
	}

	err := s.config.SetUser(cmd.Args[0])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Println("Username set successfully")
	return nil
}

// Register user command
func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("please enter a username")
	}
	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	u, _ := s.db.GetUser(context.Background(), cmd.Args[0])
	if u != nil {
		log.Fatalf("User alreadu exists")
	}

	s.db.CreateUser(context.Background(), newUser)
	return nil
}
