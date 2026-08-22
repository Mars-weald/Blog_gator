package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Mars-weald/Blog-gator/gator/internal/config"
	"github.com/Mars-weald/Blog-gator/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	gconf, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR reading: %s\n", err)
	}

	// Open database
	db, err := sql.Open("postgres", "postgres://postgres:nyanpirate@localhost:5432/gator")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dbQueries := database.New(db)

	mainState := state{
		db:   dbQueries,
		conf: &gconf,
	}

	commander := commands{
		array: map[string]func(s *state, cmd command) error{},
	}
	// register funcs needed
	commander.register("login", handlerLogin)
	commander.register("register", handlerRegister)
	commander.register("reset", handlerReset)
	commander.register("users", handlerUsers)

	//get user arguments for use
	argus := os.Args

	if len(argus) < 2 {
		fmt.Println("ERROR: too few arguments")
		os.Exit(1)
	}
	// Break user input into name and arguments for use in funcs
	comms := command{}
	comms.name = argus[1]
	comms.arguments = argus[2:]

	erroar := commander.run(&mainState, comms)
	if erroar != nil {
		fmt.Println(erroar)
		os.Exit(1)
	}
}
