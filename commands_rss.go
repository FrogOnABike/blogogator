package main

import (
	"context"
	"database/sql"
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

// *** General helper functions ***

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

// Scrape feeds - Fetches the oldest (or not yet fetched) feeds from the database, and prints the items to console
func scrapeFeeds(s *state) error {
	// Retreive the either the oldest or an unchecked feed
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("unable to retrieve: %v", err)
	}
	// Define the struct to hold the fetch time (it's a SQL NullTime type, so need to do it this way...)
	fetchedAt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	// Add that struct to the one we need to pass to the function!
	markedFeed := database.MarkFeedFetchedParams{
		ID:            nextFeed.ID,
		LastFetchedAt: fetchedAt,
	}
	// Call the function to update the fetch time for the feed
	err = s.db.MarkFeedFetched(context.Background(), markedFeed)
	if err != nil {
		return fmt.Errorf("unable to update entry: %v", err)
	}
	// Fetch the actual feed!
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fetchedFeed, err := fetchFeed(ctx, nextFeed.Url)
	if err != nil {
		return fmt.Errorf("unable to retrieve feed: %v", err)
	}
	// Print feeds item titles to console
	for _, item := range fetchedFeed.Channel.Item {
		// fmt.Printf("Item Index %v\n", i)
		fmt.Printf("Title: %s\n", item.Title)
		// fmt.Printf("Description: %s\n", item.Description)
	}
	return nil
}

// *** Handler Functions ***

// Handler for aggregating a feed - **CURRENTLY USES A STATIC FEED**
func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Please enter a duration such as 15m or 1h")
		return nil
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("unable to parse duration: %v", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// aggFeed, err := fetchFeed(ctx, "https://www.wagslane.dev/index.xml")
	// if err != nil {
	// 	log.Fatalf("Error retrieving feed: %v\n", err)
	// }
	// // fmt.Println(aggFeed)
	// // Nicely format the output so it's simpler to read and debug!
	// fmt.Printf("Title: %s\n", aggFeed.Channel.Title)
	// fmt.Printf("Description: %s\n", aggFeed.Channel.Description)
	// for i, item := range aggFeed.Channel.Item {
	// 	fmt.Printf("Item Index %v\n", i)
	// 	fmt.Printf("Title: %s\n", item.Title)
	// 	fmt.Printf("Description: %s\n", item.Description)
	// }
}

// Handler for adding a feed to the database - associates with the currently logged in user
func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		log.Fatalf("Please enter a feed name and URL")
	}
	// Create the struct to hold data for new feed to be created
	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    user.ID,
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

	// Build out the new feed_follows record
	newFF := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    createdFeed.UserID,
		FeedID:    createdFeed.ID,
	}
	_, err = s.db.CreateFeedFollow(context.Background(), newFF)
	if err != nil {
		log.Fatalf("Unable to create entry: %v\n", err)
	}
	return nil
}

// Handler for displaying all feeds in the database, along with the user who created them
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

// Handler for creating feed follows entries
func handlerFollow(s *state, cmd command, user database.User) error {
	// Check we have *something* in the args **Could add validation here for URLS?**
	if len(cmd.Args) < 1 {
		log.Fatalf("Please specify an URL to follow")
	}
	// Get feed details if it exists so can use ID
	feedDetails, err := s.db.GetFeed(context.Background(), cmd.Args[0])
	if err != nil {
		log.Fatalf("Unable to retrieve feed details: %v\n", err)
	}
	// Build out the new feed_follows record
	newFF := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedDetails.ID,
	}
	ffEntry, err := s.db.CreateFeedFollow(context.Background(), newFF)
	if err != nil {
		log.Fatalf("Unable to create entry: %v\n", err)
	}
	fmt.Printf("Feed: %s\n", ffEntry.Feedname)
	fmt.Printf("User: %s\n", ffEntry.Username)
	return nil
}

// Handler for listing all feeds current user follows
func handlerFollowing(s *state, cmd command, user database.User) error {
	userFeeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		log.Fatalf("Unable to get feeds: %v\n", err)
	}
	if len(userFeeds) == 0 {
		fmt.Println("No feeds currently followed")
		return nil
	}
	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, item := range userFeeds {
		fmt.Printf("* %s\n", item.Feedname)
	}
	return nil
}

// Handler for unfollowing a feed for the logged in user
func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		fmt.Println("Please specify a feed URL")
		return nil
	}
	ufItem := database.UnFollowFeedParams{
		Url:    cmd.Args[0],
		UserID: user.ID,
	}
	err := s.db.UnFollowFeed(context.Background(), ufItem)
	if err != nil {
		return fmt.Errorf("unable to unfollow: %v", err)
	}
	return nil
}
