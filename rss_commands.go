package main

import (
	"context"
	"encoding/xml"
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
	// Create the req
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	// Set custom User-Agent header to our app
	req.Header.Set("User-Agent", "blogogator")
	// Make the req
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()
	// Read body of response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}
	// Unmarshall the XML into the struct
	rtnRSS := RSSFeed{}

	if err := xml.Unmarshal(body, &rtnRSS); err != nil {
		return &RSSFeed{}, err
	}
	// Unescape Title and Description fields
	rtnRSS.Channel.Description = html.UnescapeString(rtnRSS.Channel.Description)
	rtnRSS.Channel.Title = html.UnescapeString(rtnRSS.Channel.Title)
	for _, item := range rtnRSS.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}
	return &rtnRSS, nil
}
