package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/CamusSisyphus/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't get next feeds to fetch", err)
		return
	}
	log.Println("Found a feed to fetch!")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			FeedID:    feed.ID,
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			Url:         item.Link,
			PublishedAt: publishedAt,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 || len(cmd.Args) > 2 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}
	time_between_reqs := cmd.Args[0]

	time_duration, err := time.ParseDuration(time_between_reqs)

	if err != nil {
		return fmt.Errorf("couldn't parse time %s: %w", time_between_reqs, err)
	}
	fmt.Printf("Collecting feeds every %s\n", time_between_reqs)

	ticker := time.NewTicker(time_duration)

	for ; ; <-ticker.C {
		scrapeFeeds(s)

	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	if len(cmd.Args) < 2 {
		return errors.New("Requires Feed name and url!")
	}
	feed_name := cmd.Args[0]
	feed_url := cmd.Args[1]

	rssFeed, err := s.db.CreateFeed(context.Background(),
		database.CreateFeedParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Name:      feed_name,
			Url:       feed_url,
			UserID:    user.ID,
		})

	if err != nil {
		return fmt.Errorf("couldn't create feed entry: %w", err)
	}

	_, err = s.db.CreateFeedFollow(context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			UserID:    user.ID,
			FeedID:    rssFeed.ID,
		})
	if err != nil {
		return fmt.Errorf("couldn't create feed_follow entry: %w", err)
	}

	fmt.Println(rssFeed)
	return nil
}

func handlerFeeds(s *state, cmd command) error {

	if len(cmd.Args) > 0 {
		return errors.New("No args required!")
	}

	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Cannot retrieve feeds: %w", err)
	}

	fmt.Println("Feed List:")
	for _, feed := range feeds {
		user, err := s.db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("Cannot retrieve user by ID: %w", err)
		}
		fmt.Printf("	Feed Name: %s\n", feed.Name)
		fmt.Printf("		- url: %s\n", feed.Url)
		fmt.Printf("		- user: %s\n", user.Name)
	}

	return nil
}
