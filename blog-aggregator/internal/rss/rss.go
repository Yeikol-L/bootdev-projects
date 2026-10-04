package rss

import (
	"context"
	"encoding/xml"
	"html"
	"net/http"
)

type Feed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Item        []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*Feed, error) {
	var feed *Feed
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return feed, nil
	}
	req.Header.Set("User-Agent", "gator")
	req.Header.Set("Accept-Content", "application/rss+xml")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return feed, nil
	}
	defer res.Body.Close()
	err = xml.NewDecoder(res.Body).Decode(&feed)

	if err != nil {
		return feed, nil
	}
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return feed, nil
}
