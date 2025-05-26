package main

import (
	"io"
	"fmt"
	"html"
	"context"
	"net/http"
	"encoding/xml"
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
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return &RSSFeed{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if resp.StatusCode > 299 {
		return &RSSFeed{}, fmt.Errorf("Request returned with status code: %d\n", resp.StatusCode)
	}
	if err != nil {
		return &RSSFeed{}, err
	}

	var rss_feed RSSFeed
	err = xml.Unmarshal(body, &rss_feed)
	if err != nil {
		return &RSSFeed{}, err
	}

	rss_feed.Channel.Title = html.UnescapeString(rss_feed.Channel.Title)
	rss_feed.Channel.Description = html.UnescapeString(rss_feed.Channel.Description)
	for i, item := range rss_feed.Channel.Item {
		rss_feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		rss_feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return &rss_feed, nil
}
