package json

import (
	"encoding/json"

	"github.com/jadmadi/Go-OPML/internal/opml"
	"github.com/jadmadi/Go-OPML/pkg/rss"
)

// PodcastFeed represents a simplified structure for JSON output
type PodcastFeed struct {
	Title    string        `json:"title"`
	URL      string        `json:"url"`
	Episodes []rss.Episode `json:"episodes,omitempty"`
}

// GenerateJSON converts the OPML struct to a slice of PodcastFeed and marshals it to JSON
func GenerateJSON(opmlData *opml.OPML, fetchRSS bool, maxEpisodes int, timeout string, concurrency int) ([]byte, error) {
	var feeds []PodcastFeed

	for _, outline := range opmlData.Body.Outlines {
		feed := PodcastFeed{
			Title: outline.Title,
			URL:   outline.XMLUrl,
		}

		// Fetch RSS episodes if requested
		if fetchRSS && outline.XMLUrl != "" {
			episodes, err := rss.FetchEpisodes(outline.XMLUrl, maxEpisodes, timeout)
			if err == nil {
				feed.Episodes = episodes
			}
		}

		feeds = append(feeds, feed)
	}

	return json.MarshalIndent(feeds, "", "  ")
}
