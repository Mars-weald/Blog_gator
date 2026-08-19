package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Mars-weald/Blog-gator/gator/internal/config"
	"github.com/Mars-weald/Blog-gator/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db   *database.Queries
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

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("ERROR: no argument to register")
	}

	panams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
	}
	//Check if user exists in database
	x, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println("ERROR checking databse during registration")
		os.Exit(1)
	}
	if x.Name == panams.Name {
		fmt.Println("ERROR: User already registered")
		os.Exit(1)
	}

	_, err = s.db.CreateUser(context.Background(), panams)
	if err != nil {
		fmt.Println("ERROR creating user during registration")
		os.Exit(1)
	}
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
