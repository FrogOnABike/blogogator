package main

import (
	"context"
	"fmt"
	"log"
)

// Reset command
func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		log.Fatalf("error resetting database")
	}
	fmt.Println("Users databse reset")
	return nil
}
