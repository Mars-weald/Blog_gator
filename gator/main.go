package main

import (
	"fmt"

	"github.com/Mars-weald/Blog-gator/gator/internal/config"
)

func main() {
	gconf, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR reading: %s\n", err)
	}

	gconf.SetUser("Mars")

	data, err := config.Read()
	if err != nil {
		fmt.Printf("ERROR reading: %s\n", err)
	}

	fmt.Printf("%+v", data)
}
