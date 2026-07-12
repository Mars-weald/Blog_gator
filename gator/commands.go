package main

import (
	"fmt"

	"github.com/Mars-weald/Blog-gator/gator/internal/config"
)

type state struct {
	conf *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	array map[string]func(*state, command) error
}

func handlerLogin(s *state, cmd command) error {
	if cmd.arguments == nil {
		return fmt.Errorf("ERROR: No argument for login")
	}

	err := s.conf.SetUser(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("ERROR logging in: set user err: %w", err)
	}
	fmt.Println("User has been set")
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	return nil
}

func (c *commands) register(s *state, cmd command) error {
	return nil
}
