package main

import (
	"fmt"
	"log"
	"os"

	"github.com/CamusSisyphus/blog_aggregator/internal/config"
)

type state struct {
	cfg *(config.Config)
}

func main() {

	c, err := config.Read()

	if err != nil {
		fmt.Println(err)
	}
	s := &state{cfg: &c}

	commands := commands{make(map[string]func(*state, command) error)}
	commands.register("login", handlerLogin)
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Not Enough Arguments Providedbootdev run dca1352a-7600-4d1d-bfdf-f9d741282e55")
	}
	commandName := args[1]
	commandArgs := args[2:]

	err = commands.run(s, command{Name: commandName, Args: commandArgs})
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
	// handlerLogin(&s, command{name: "handlerLogin", args: []string{"Mason"}})

}
