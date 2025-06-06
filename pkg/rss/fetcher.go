package rss

import (
	"time"

	"github.com/mmcdole/gofeed"
)

// Episode represents a podcast episode
type Episode struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	PublishDate time.Time `json:"publishDate"`
	Description string    `json:"description"`
}

// FetchEpisodes fetches episodes from an RSS feed
func FetchEpisodes(feedURL string, maxEpisodes int, timeout string) ([]Episode, error) {
	fp := gofeed.NewParser()

	// Parse timeout duration
	timeoutDuration, err := time.ParseDuration(timeout)
	if err != nil {
		timeoutDuration = 30 * time.Second
	}

	// Set timeout for HTTP client
	fp.Client.Timeout = timeoutDuration

	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	var episodes []Episode
	limit := maxEpisodes
	if len(feed.Items) < limit {
		limit = len(feed.Items)
	}

	for i := 0; i < limit; i++ {
		item := feed.Items[i]
		episode := Episode{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
		}

		if item.PublishedParsed != nil {
			episode.PublishDate = *item.PublishedParsed
		}

		episodes = append(episodes, episode)
	}

	return episodes, nil
}
