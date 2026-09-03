package main

import (
	"context"
	"database/sql"
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
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("ERROR: need time argument")
	}

	time_between_reqs, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("ERROR parsing duration: %w\n", err)
	}

	fmt.Printf("Aggregating feeds every %s...\n", cmd.arguments[0])

	ticker := time.NewTicker(time_between_reqs)
	for ; ; <-ticker.C {
		feedScraper(s)
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("ERROR: Need feed name and URL")
	} else if len(cmd.arguments) == 1 {
		return fmt.Errorf("ERROR: Need URL")
	}

	parms := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	}

	food, err := s.db.CreateFeed(context.Background(), parms)
	if err != nil {
		return fmt.Errorf("ERROR creating feed: %w\n", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    parms.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("ERROR creating feed-follow: %w\n", err)
	}
	fmt.Println("Follow-feed record created")
	fmt.Printf("%+v\n", food)
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feedArray, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("ERROR getting feeds: %w\n", err)
	}
	for _, feed := range feedArray {
		fmt.Printf("%s (%s) created by %s\n", feed.Name, feed.Url, feed.Name_2)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	//Get feed ID
	food, err := s.db.UrlSearch(context.Background(), cmd.arguments[0])
	if err != nil {
		return fmt.Errorf("ERROR getting feed: %w\n", err)
	}

	panams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    food.ID,
	}

	record, err := s.db.CreateFeedFollow(context.Background(), panams)
	if err != nil {
		return fmt.Errorf("ERROR making feed-follow: %w\n", err)
	}
	for _, feed := range record {
		fmt.Printf("%s (%s)\n", feed.UserName, feed.FeedNme)
	}
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedArray, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("ERROR getting feeds: %w\n", err)
	}
	for _, food := range feedArray {
		fmt.Println(food.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	url := cmd.arguments[0]
	feedDat, err := s.db.UrlSearch(context.Background(), url)
	if err != nil {
		return fmt.Errorf("ERROR getting feed: %w\n", err)
	}
	parms := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feedDat.ID,
	}
	err = s.db.Unfollow(context.Background(), parms)
	if err != nil {
		return fmt.Errorf("ERROR unfollowing: %w\n", err)
	}
	return nil
}

func feedScraper(s *state) error {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("ERROR fetching feeds: %w\n", err)
	}

	thyme := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	parms := database.MarkFeedFetchedParams{
		LastFetchedAt: thyme,
		ID:            nextFeed.ID,
	}

	err = s.db.MarkFeedFetched(context.Background(), parms)
	if err != nil {
		return fmt.Errorf("ERROR marking feed as fetched: %w\n", err)
	}

	russ, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("ERROR fetching feeds: %w\n", err)
	}
	fmt.Println(russ.Channel.Title)

	for _, item := range russ.Channel.Item {
		fmt.Printf("-- %s\n", item.Title)
	}
	return nil
}

// middleware
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
		if err != nil {
			return fmt.Errorf("ERROR getting user: %w\n", err)
		}
		return handler(s, cmd, user)
	}
}

// commands funcs
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
