package parser_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	main "github.com/jadmadi/Go-OPML/cmd/go-opml"
)

func TestParseOPML(t *testing.T) {
	// Create a temporary OPML file for testing
	tempFile, err := os.CreateTemp("", "test_opml_*.xml")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write test OPML content to the temporary file
	testOPML := `<?xml version="1.0" encoding="UTF-8"?>
<opml>
    <head/>
    <body version="1.0">
        <outline xmlUrl="https://example.com/feed1.rss" title="Test Feed 1"/>
        <outline xmlUrl="https://example.com/feed2.rss" title="Test Feed 2"/>
    </body>
</opml>`
	err = os.WriteFile(tempFile.Name(), []byte(testOPML), 0644)
	if err != nil {
		t.Fatalf("Error writing to temporary file: %v", err)
	}

	// Parse the test OPML file
	opml, err := main.ParseOPML(tempFile.Name())
	if err != nil {
		t.Fatalf("Error parsing OPML: %v", err)
	}

	// Check the parsed content
	if len(opml.Body.Outlines) != 2 {
		t.Errorf("Expected 2 outlines, got %d", len(opml.Body.Outlines))
	}

	if opml.Body.Outlines[0].Title != "Test Feed 1" {
		t.Errorf("Expected title 'Test Feed 1', got '%s'", opml.Body.Outlines[0].Title)
	}

	if opml.Body.Outlines[1].XMLUrl != "https://example.com/feed2.rss" {
		t.Errorf("Expected URL 'https://example.com/feed2.rss', got '%s'", opml.Body.Outlines[1].XMLUrl)
	}
}

func TestGenerateJSON(t *testing.T) {
	// Create a test OPML struct
	testOPML := &main.OPML{
		Body: main.Body{
			Outlines: []main.Outline{
				{XMLUrl: "https://example.com/feed1.rss", Title: "Test Feed 1"},
				{XMLUrl: "https://example.com/feed2.rss", Title: "Test Feed 2"},
			},
		},
	}

	// Generate JSON
	jsonData, err := main.GenerateJSON(testOPML)
	if err != nil {
		t.Fatalf("Error generating JSON: %v", err)
	}

	// Check the generated JSON
	expectedJSON := `[
  {
    "title": "Test Feed 1",
    "url": "https://example.com/feed1.rss"
  },
  {
    "title": "Test Feed 2",
    "url": "https://example.com/feed2.rss"
  }
]`

	if string(jsonData) != expectedJSON {
		t.Errorf("Generated JSON does not match expected output.\nExpected: %s\nGot: %s", expectedJSON, string(jsonData))
	}
}

func TestFetchRSSFeed(t *testing.T) {
	// Test with a timeout that should work for most cases
	timeout := 10 * time.Second
	maxEpisodes := 5

	// Test with a well-known RSS feed (this is an integration test)
	// Note: This test requires internet connection
	episodes, err := main.FetchRSSFeed("https://feeds.simplecast.com/e_GRxR9a", maxEpisodes, timeout)

	// If the feed is accessible, we should get episodes
	if err == nil {
		if len(episodes) > maxEpisodes {
			t.Errorf("Expected maximum %d episodes, got %d", maxEpisodes, len(episodes))
		}

		// Check if episodes have required fields
		for i, episode := range episodes {
			if episode.Title == "" {
				t.Errorf("Episode %d has empty title", i)
			}
			if episode.PublishDate == "" {
				t.Errorf("Episode %d has empty publish date", i)
			}
		}
	} else {
		t.Logf("RSS feed test skipped due to network error: %v", err)
	}
}

func TestPodcastStructures(t *testing.T) {
	// Test JSON marshaling of podcast structures
	podcast := main.PodcastWithEpisodes{
		Title: "Test Podcast",
		URL:   "https://example.com/feed.rss",
		Episodes: []main.EpisodeInfo{
			{
				Title:       "Episode 1",
				Link:        "https://example.com/episode1",
				PublishDate: "2024-01-01T12:00:00Z",
				Description: "Test episode description",
			},
		},
	}

	jsonData, err := json.Marshal(podcast)
	if err != nil {
		t.Fatalf("Error marshaling podcast to JSON: %v", err)
	}

	var unmarshaled main.PodcastWithEpisodes
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Error unmarshaling JSON: %v", err)
	}

	if unmarshaled.Title != podcast.Title {
		t.Errorf("Expected title '%s', got '%s'", podcast.Title, unmarshaled.Title)
	}

	if len(unmarshaled.Episodes) != 1 {
		t.Errorf("Expected 1 episode, got %d", len(unmarshaled.Episodes))
	}
}
