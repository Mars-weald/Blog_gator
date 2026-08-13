package main

import (
	"fmt"
	"os"

	"github.com/Mars-weald/Blog-gator/gator/internal/config"
)

func main() {
	gconf, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR reading: %s\n", err)
	}

	mainState := state{
		conf: &gconf,
	}

	commander := commands{
		array: map[string]func(s *state, cmd command) error{},
	}

	commander.register("login", handlerLogin)

	argus := os.Args

	if len(argus) < 2 {
		fmt.Println("ERROR: too few arguments")
		os.Exit(1)
	}

	comms := command{}
	comms.name = argus[1]
	comms.arguments = argus[2:]

	erroar := commander.run(&mainState, comms)
	if erroar != nil {
		fmt.Println(erroar)
		os.Exit(1)
	}
}
