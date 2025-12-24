package main

import (
	"net/http"
	"time"

	"encoding/xml"
)

type RSSFeed struct {
	Channel struct {
		Title         string    `xml:"title"`
		Link          string    `xml:"link"`
		Description   string    `xml:"description"`
		Language      string    `xml:"language"`
		Generator     string    `xml:"generator"`
		LastBuildDate string    `xml:"lastBuildDate"`
		Items         []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}
	defer resp.Body.Close()

	feed := RSSFeed{}
	err = xml.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return RSSFeed{}, err
	}
	return feed, nil
}
