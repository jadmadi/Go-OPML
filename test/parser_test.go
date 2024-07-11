package main

import (
	"os"
	"testing"
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
	opml, err := ParseOPML(tempFile.Name())
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
	testOPML := &OPML{
		Body: Body{
			Outlines: []Outline{
				{XMLUrl: "https://example.com/feed1.rss", Title: "Test Feed 1"},
				{XMLUrl: "https://example.com/feed2.rss", Title: "Test Feed 2"},
			},
		},
	}

	// Generate JSON
	jsonData, err := GenerateJSON(testOPML)
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
