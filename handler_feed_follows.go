package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CamusSisyphus/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {

	if len(cmd.Args) < 1 {
		return errors.New("Requires Feed url!")
	}

	feed_url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("No feed exists with url: %s", feed_url)
	}

	_, err = s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		})

	if err != nil {
		return fmt.Errorf("couldn't create feed_follow entry: %w", err)
	}

	fmt.Printf("Feed: %s (%s) has been followed by %s", feed.Name, feed_url, user.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	if len(cmd.Args) > 0 {
		return errors.New("No args required!")
	}

	feed_follows_for_user, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("Failed to retreive follow: %w", err)
	}

	fmt.Println("Current following feeds:")
	for _, feed_follow := range feed_follows_for_user {
		fmt.Printf("	- %s", feed_follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {

	if len(cmd.Args) < 1 {
		return errors.New("Requires Feed url!")
	}

	feed_url := cmd.Args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), feed_url)
	if err != nil {
		return fmt.Errorf("No feed exists with url: %s", feed_url)
	}

	err = s.db.DeleteFeedFollow(context.Background(),
		database.DeleteFeedFollowParams{
			UserID: user.ID,
			FeedID: feed.ID,
		})
	if err != nil {
		return fmt.Errorf("couldn't delete feed_follow entry: %w", err)
	}

	fmt.Printf("Feed: %s (%s) has been unfollowed by %s", feed.Name, feed_url, user.Name)
	return nil
}
