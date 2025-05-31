package main

import (
	"io"
	"fmt"
	"html"
	"time"
	"context"
	"net/http"
	"database/sql"
	"encoding/xml"
	"github.com/t6kke/gator/internal/database"
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

func scrapeFeeds(s *state) error {
	new_ctx := context.Background()
	next_feed, err := s.dbq.GetNextFeedToFetch(new_ctx)
	if err != nil {
		return err
	}

	feed_url := next_feed.Url
	feed_id := next_feed.ID
	current_time := time.Now()
	nullable_time := sql.NullTime{
		Time:  current_time,
		Valid: true,
	}
	marking_params := database.MarkFeedFetchedParams{
		UpdatedAt:     current_time,
		LastFetchedAt: nullable_time,
		ID:            feed_id,
	}
	_, err = s.dbq.MarkFeedFetched(new_ctx, marking_params)
	if err != nil {
		return err
	}

	feed, err := fetchFeed(new_ctx, feed_url)
	if err != nil {
		return err
	}

	fmt.Println(feed.Channel.Link)
	fmt.Println(feed.Channel.Title)
	fmt.Println(feed.Channel.Description)
	fmt.Println("-------------------------------------------------------")
	fmt.Println("items:")
	for i, item := range feed.Channel.Item {
		fmt.Printf("Item: %d --- Link: %s\n", i+1, item.Link)
		fmt.Printf("Title: '%s'\n", item.Title)
		//fmt.Println(item.Description)
		fmt.Println("-------------------------------------------------------")
	}

	return nil
}
