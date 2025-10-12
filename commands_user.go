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

// Login command
func handlerLogin(s *state, cmd command) error {
	// Check that we only have a single username in the args slice, otherwise return an error
	if len(cmd.Args) != 1 {
		log.Fatalf("please enter a username")
	}

	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		log.Fatalf("User not found")
	}
	err = s.config.SetUser(cmd.Args[0])
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
	_, err := s.db.GetUser(context.Background(), cmd.Args[0])
	// If the GetUser func finds a match then a nil error is returned
	// So nil err indicates an existing (duplicate) username
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	createdUser, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return fmt.Errorf("unable to create user: %v", err)
	}
	err = s.config.SetUser(createdUser.Name)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	fmt.Printf("User created successfully: %s\n", createdUser.Name)
	return nil
}
