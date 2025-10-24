package gogator

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
	"time"
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
