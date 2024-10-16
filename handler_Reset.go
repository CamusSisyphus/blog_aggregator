package main

import (
	"context"
	"errors"
	"fmt"
)

func handlerReset(s *state, cmd command) error {

	err := s.db.Reset(context.Background())
	if err != nil {
		return errors.New("Failed to reset data base!")
	}

	fmt.Println("User table has been reset")
	return nil
}
