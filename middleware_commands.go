package main

import (
	"context"
	"errors"

	"github.com/frogonabike/blogogator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, u database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		name := s.config.CurrentUserName
		if name == "" {
			return errors.New("not logged in")
		}
		u, err := s.db.GetUser(context.Background(), name)
		if err != nil {
			return err
		}
		return handler(s, cmd, u)
	}
}
