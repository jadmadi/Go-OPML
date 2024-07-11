package main

import (
	"encoding/json"
)

// PodcastFeed represents a simplified structure for JSON output
type PodcastFeed struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// GenerateJSON converts the OPML struct to a slice of PodcastFeed and marshals it to JSON
func GenerateJSON(opml *OPML) ([]byte, error) {
	var feeds []PodcastFeed

	for _, outline := range opml.Body.Outlines {
		feed := PodcastFeed{
			Title: outline.Title,
			URL:   outline.XMLUrl,
		}
		feeds = append(feeds, feed)
	}

	return json.MarshalIndent(feeds, "", "  ")
}
