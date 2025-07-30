package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mmnsrti/rssproject/internal/database"
)

func startScraping(
	db *database.Queries,
	concurrency int,
	timeBetweenRequest time.Duration,
) {
	log.Printf("Starting scraping with concurrency: %d and time between requests: %s", concurrency, timeBetweenRequest)
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		log.Println("Scraping feeds...")
		feed, err := db.GetFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("Error fetching feeds: %v", err)
			continue
		}
		// Process the feeds
		wg := &sync.WaitGroup{}
		for _, f := range feed {
			wg.Add(1)
			go scrapeFeed(wg, db, f)

		}
		wg.Wait()
	}

}
func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()
	// Implement the logic to scrape the feed
	log.Println("Scraping feed...")
	// Simulate some work
	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
		return
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching RSS feed: %v", err)
		return
	}
	for _, item := range rssFeed.Channel.Items {
		log.Printf("Processing item: %s", item.Title)
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing publication date: %v", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			Title:       item.Title,
			Description: description.String, // Use description.String if CreatePostParams expects string
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				log.Printf("Post already exists: %s", item.Title)
				continue
			}
			log.Printf("Error creating post: %v", err)

		}

	}

	log.Println("Feed scraped successfully")
}
