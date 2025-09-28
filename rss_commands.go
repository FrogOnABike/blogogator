package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/frogonabike/blogogator/internal/database"
	"github.com/google/uuid"
)

// Structs to unmarshal XML data into
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

// Fetch an RSS feed from a given URL and return it as a struct
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
	for i := range rtnRSS.Channel.Item {
		rtnRSS.Channel.Item[i].Title = html.UnescapeString(rtnRSS.Channel.Item[i].Title)
		rtnRSS.Channel.Item[i].Description = html.UnescapeString(rtnRSS.Channel.Item[i].Description)
	}
	return &rtnRSS, nil
}

// Handler for aggregating a feed - **CURRENTLY USES A STATIC FEED**
func handlerAgg(s *state, cmd command) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	aggFeed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatalf("Error retrieving feed: %v\n", err)
	}
	// fmt.Println(aggFeed)
	// Nicely format the output so it's simpler to read and debug!
	fmt.Printf("Title: %s\n", aggFeed.Channel.Title)
	fmt.Printf("Description: %s\n", aggFeed.Channel.Description)
	for i, item := range aggFeed.Channel.Item {
		fmt.Printf("Item Index %v\n", i)
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Description: %s\n", item.Description)
	}
	return nil
}

// Handler for adding a feed to the database - associates with the currently logged in user
func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		log.Fatalf("Please enter a feed name and URL")
	}
	// Obtain details for the currently logged in user as per config file
	curUser, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		log.Fatalf("Unable to locate user: %v\n", err)
	}
	// Create the struct to hold data for new feed to be created
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    curUser.ID,
	}
	// Insert the new feed into db!
	createdFeed, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		log.Fatalf("Unable to add feed: %v\n", err)
	}
	fmt.Printf("ID: %v\n", createdFeed.ID)
	fmt.Printf("Created at: %v\n", createdFeed.CreatedAt)
	fmt.Printf("Updated at: %v\n", createdFeed.UpdatedAt)
	fmt.Printf("Name: %v\n", createdFeed.Name)
	fmt.Printf("URL: %v\n", createdFeed.Url)
	fmt.Printf("User ID: %v\n", createdFeed.UserID)

	return nil
}

// Handle for displaying all feeds in the database, along with the user who created them
func handlerFeeds(s *state, cmd command) error {
	feedsList, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalf("Error retrieving feeds: %v\n", err)
	}
	for _, feed := range feedsList {
		fmt.Printf("Feed: %s URL: %s Created By: %s\n", feed.Feedname, feed.Url, feed.Username)
	}
	return nil
}
