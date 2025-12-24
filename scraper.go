package main

import (
	"log"
	"time"

	"github.com/keithfy96/go-project/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {

	log.Printf("Scrapping on %v goroutines every %s duration", concurrency, timeBetweenRequests)
}
