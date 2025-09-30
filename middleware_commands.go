package main

import "github.com/frogonabike/blogogator/internal/database"

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

}
