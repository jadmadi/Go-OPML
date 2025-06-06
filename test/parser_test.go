package test

import (
	"testing"

	"github.com/jadmadi/Go-OPML/internal/json"
	"github.com/jadmadi/Go-OPML/internal/opml"
	"github.com/jadmadi/Go-OPML/pkg/rss"
)

func TestParseOPML(t *testing.T) {
	// Test parsing OPML file
	opmlData, err := opml.ParseOPML("../examples/sample.opml")
	if err != nil {
		t.Fatalf("Failed to parse OPML: %v", err)
	}

	if len(opmlData.Body.Outlines) == 0 {
		t.Error("No outlines found in OPML")
	}
}

func TestGenerateJSON(t *testing.T) {
	// Test JSON generation
	opmlData, err := opml.ParseOPML("../examples/sample.opml")
	if err != nil {
		t.Fatalf("Failed to parse OPML: %v", err)
	}

	jsonData, err := json.GenerateJSON(opmlData, false, 10, "30s", 5)
	if err != nil {
		t.Errorf("Failed to generate JSON: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("Generated JSON is empty")
	}
}

func TestFetchRSS(t *testing.T) {
	// Test RSS fetching (skip if no internet connection)
	testURL := "https://feeds.megaphone.fm/thediaryofaceo"
	episodes, err := rss.FetchEpisodes(testURL, 5, "30s")
	if err != nil {
		t.Skipf("Skipping RSS test due to error: %v", err)
		return
	}

	if len(episodes) == 0 {
		t.Log("No episodes fetched (this might be normal)")
	}
}
