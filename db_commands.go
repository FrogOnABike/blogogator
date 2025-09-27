package main

import (
	"context"
	"fmt"
	"log"
)

// Reset command - Clears out the user database! **USE WITH CAUTION IN PROD!**
func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		log.Fatalf("error resetting database: %v\n", err)
	}
	fmt.Println("Users databse reset")
	return nil
}

// GetUsers command - List the configured users and indicate who is currently logged in
func handlerGetUsers(s *state, cmd command) error {
	userList, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("error retrieving users: %v\n", err)
	}
	for _, user := range userList {
		if user.Name == s.config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}
