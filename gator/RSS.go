package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	feed := RSSFeed{}
	client := http.DefaultClient

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &feed, fmt.Errorf("ERROR making request: %w", err)
	}

	req.Header.Set("User-Agent", "gator")

	res, err := client.Do(req)
	if err != nil {
		return &feed, fmt.Errorf("ERROR getting response: %w", err)
	}

	bites, err := io.ReadAll(res.Body)
	if err != nil {
		return &feed, fmt.Errorf("ERROR reading body: %w", err)
	}

	err = xml.Unmarshal(bites, &feed)
	if err != nil {
		return &feed, fmt.Errorf("ERROR unmarshalling: %w", err)
	}
	//Decoding time!
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	feed.Channel.Item[0].Title = html.UnescapeString(feed.Channel.Item[0].Title)
	feed.Channel.Item[0].Description = html.UnescapeString(feed.Channel.Item[0].Description)

	return &feed, nil
}
