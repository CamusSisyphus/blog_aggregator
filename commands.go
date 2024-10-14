package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	availableCommands map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.availableCommands[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.availableCommands[cmd.Name]
	if ok {
		return f(s, cmd)
	}
	return errors.New("No such Command!")
}
