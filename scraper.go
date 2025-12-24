package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/keithfy96/go-project/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {

	log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Error getting next feeds to fetch:", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(wg, db, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *database.Queries, feed database.Feed) {
	defer wg.Done()
	log.Printf("Scraping feed: %s", feed.Name)

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}
	RSSFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error getting feed:", err, "on feed", feed.Name)
		return
	}
	for _, item := range RSSFeed.Channel.Items {
		t, err := time.Parse(time.RFC1123, item.PubDate)
		if err != nil {
			log.Printf("Error parsing published date for item %s: %v", item.Title, err)
			return
		}
		// log.Printf("item %s collected, %s", item.Title, item.Link)
		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Description: sql.NullString{
				String: item.Description,
				Valid:  item.Description != "",
			},
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Error creating post:", err)
			return
		}
	}
	log.Printf("feed %s collected, %v posts found", feed.Name, len(RSSFeed.Channel.Items))
}
