package main

import (
	"fmt"
	"os"

	"github.com/Mars-weald/Blog-gator/gator/internal/config"
)

const configFile = "/home/mars_orleis/.gatorconfig.json"

func main() {
	gconf, err := config.Read(configFile)
	if err != nil {
		return
	}

	gconf.SetUser("Mars")

	data, err := os.ReadFile(configFile)
	if err != nil {
		return
	}
	fmt.Print(string(data))
}
