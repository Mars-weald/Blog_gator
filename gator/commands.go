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
	if len(cmd.arguments) == 0 {
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
	funcToRun, ok := c.array[cmd.name]
	if !ok {
		return fmt.Errorf("Command not found")
	} else {
		err := funcToRun(s, cmd)
		if err != nil {
			return fmt.Errorf("ERROR running: %w", err)
		}
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	_, ok := c.array[name]
	if ok {
		fmt.Println("Function already registered")
	} else {
		c.array[name] = f
	}
}
