package gogator

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/prchop/gogator/internal/database"
)

type RSSFeed struct {
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title string    `xml:"title"`
	Link  string    `xml:"link"`
	Desc  string    `xml:"description"`
	Item  []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Desc    string `xml:"description"`
	PubDate string `xml:"pubDate"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		html.EscapeString(feedURL),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "gogator")

	client := http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err := xml.Unmarshal(b, &feed); err != nil {
		return nil, err
	}

	fc := feed.Channel
	fc.Title = html.UnescapeString(fc.Title)
	fc.Desc = html.UnescapeString(fc.Desc)

	for i := range fc.Item {
		fc.Item[i].Title = html.UnescapeString(fc.Item[i].Title)
		fc.Item[i].Desc = html.UnescapeString(fc.Item[i].Desc)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) {
	ctx := context.Background()
	dbFeed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		log.Printf("couldn't get feed: %v\n", err)
		return
	}

	err = s.db.MarkFeedFetched(ctx,
		database.MarkFeedFetchedParams{
			UpdatedAt: time.Now().UTC(),
			UserID:    dbFeed.UserID,
		},
	)
	if err != nil {
		log.Printf("couldn't mark feed to fetched: %v\n", err)
		return
	}

	rss, err := fetchFeed(ctx, dbFeed.Url)
	if err != nil {
		log.Printf("couldn't fetch feed: %v\n", err)
		return
	}

	for _, ri := range rss.Channel.Item {
		fmt.Printf("RSS Title: %s\n", ri.Title)
	}
}
