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

	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println("ERROR: user not registered")
		os.Exit(1)
	}

	err = s.conf.SetUser(cmd.arguments[0])
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
	_, err := s.db.GetUser(context.Background(), cmd.arguments[0])
	if err == nil {
		fmt.Println("ERROR: User already registered")
	}

	_, err = s.db.CreateUser(context.Background(), panams)
	if err != nil {
		return fmt.Errorf("ERROR creating user during registry: %w", err)
	}

	err = s.conf.SetUser(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("ERROR registering: set user err: %w", err)
	}
	fmt.Println("User registered")
	fmt.Println("User has been set")
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	if err != nil {
		fmt.Println("ERROR resetting: %w", err)
		os.Exit(1)
	}
	fmt.Println("Database successfully reset")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	names, err := s.db.GetUsers(context.Background())
	if err != nil {
		fmt.Println("ERROR getting data from database")
		os.Exit(1)
	}
	for i := 0; i < len(names); i++ {
		if names[i] == s.conf.CurrentUserName {
			fmt.Printf("* %s (current)\n", names[i])
		} else {
			fmt.Println(names[i])
		}
	}
	return nil
}

func handlerAggregate(s *state, cmd command) error {
	reallySimpleFeed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("ERROR: %w", err)
	}
	fmt.Printf("%+v\n", reallySimpleFeed)
	return nil
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("ERROR: Need feed name and URL")
	} else if len(cmd.arguments) == 1 {
		return fmt.Errorf("ERROR: Need URL")
	}

	user := s.conf.CurrentUserName
	x, err := s.db.GetUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("ERROR getting user: %w\n", err)
	}

	parms := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    x.ID,
	}

	food, err := s.db.CreateFeed(context.Background(), parms)
	if err != nil {
		return fmt.Errorf("ERROR creating feed: %w\n", err)
	}
	fmt.Printf("%+v\n", food)
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
