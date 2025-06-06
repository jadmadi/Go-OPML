package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mmcdole/gofeed"
)

type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Version string   `xml:"version,attr"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

type Head struct {
	Title string `xml:"title,omitempty"`
}

type Body struct {
	Outlines []Outline `xml:"outline"`
}

type Outline struct {
	XMLUrl string `xml:"xmlUrl,attr"`
	Title  string `xml:"title,attr"`
}

type PodcastWithEpisodes struct {
	Title    string        `json:"title"`
	URL      string        `json:"url"`
	Episodes []EpisodeInfo `json:"episodes,omitempty"`
}

type EpisodeInfo struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	PublishDate string `json:"publishDate"`
	Description string `json:"description"`
}

func main() {
	// Define flags
	inputFile := flag.String("input", "", "Path to the input OPML file")
	outputFile := flag.String("output", "output.json", "Path to the output JSON file")
	fetchRSS := flag.Bool("fetch-rss", false, "Fetch RSS feeds and include episode information")
	maxEpisodes := flag.Int("max-episodes", 10, "Maximum number of episodes to fetch per podcast")
	timeout := flag.Duration("timeout", 30*time.Second, "Timeout for fetching RSS feeds")
	concurrency := flag.Int("concurrency", 5, "Number of concurrent RSS feed fetches")

	// Custom usage function
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Go-OPML: OPML to JSON converter with RSS feed fetching capabilities\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  Convert OPML to JSON:\n")
		fmt.Fprintf(os.Stderr, "    %s -input podcasts.opml -output podcasts.json\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Convert OPML to JSON and fetch RSS feeds:\n")
		fmt.Fprintf(os.Stderr, "    %s -input podcasts.opml -output podcasts.json -fetch-rss\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  Fetch RSS feeds with custom settings:\n")
		fmt.Fprintf(os.Stderr, "    %s -input podcasts.opml -output podcasts.json -fetch-rss -max-episodes 20 -timeout 1m -concurrency 10\n", os.Args[0])
	}

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: Input OPML file is required")
		flag.Usage()
		os.Exit(1)
	}

	// Parse the OPML file
	opml, err := ParseOPML(*inputFile)
	if err != nil {
		log.Fatalf("Error parsing OPML file: %v", err)
	}

	var podcasts []PodcastWithEpisodes
	var wg sync.WaitGroup
	podcastChan := make(chan PodcastWithEpisodes, len(opml.Body.Outlines))
	semaphore := make(chan struct{}, *concurrency)

	for _, outline := range opml.Body.Outlines {
		wg.Add(1)
		go func(outline Outline) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			podcast := PodcastWithEpisodes{
				Title: outline.Title,
				URL:   outline.XMLUrl,
			}

			if *fetchRSS {
				episodes, err := FetchRSSFeed(outline.XMLUrl, *maxEpisodes, *timeout)
				if err != nil {
					log.Printf("Error fetching RSS for %s: %v", outline.Title, err)
				} else {
					podcast.Episodes = episodes
				}
			}

			podcastChan <- podcast
		}(outline)
	}

	go func() {
		wg.Wait()
		close(podcastChan)
	}()

	for podcast := range podcastChan {
		podcasts = append(podcasts, podcast)
	}

	// Generate JSON
	jsonData, err := json.MarshalIndent(podcasts, "", "  ")
	if err != nil {
		log.Fatalf("Error generating JSON: %v", err)
	}

	// Write JSON to file
	err = os.WriteFile(*outputFile, jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON to file: %v", err)
	}

	fmt.Printf("Successfully processed OPML and generated JSON. Output written to %s\n", *outputFile)
}

// Export ParseOPML function for testing
func ParseOPML(filename string) (*OPML, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var opml OPML
	err = xml.Unmarshal(data, &opml)
	if err != nil {
		return nil, err
	}

	return &opml, nil
}

// Export FetchRSSFeed function for testing
func FetchRSSFeed(feedURL string, maxEpisodes int, timeout time.Duration) ([]EpisodeInfo, error) {
	fp := gofeed.NewParser()
	fp.Client = &http.Client{Timeout: timeout}
	feed, err := fp.ParseURL(feedURL)
	if err != nil {
		return nil, err
	}

	var episodes []EpisodeInfo
	for i, item := range feed.Items {
		if i >= maxEpisodes {
			break
		}
		episode := EpisodeInfo{
			Title:       item.Title,
			Link:        item.Link,
			PublishDate: item.PublishedParsed.Format(time.RFC3339),
			Description: item.Description,
		}
		episodes = append(episodes, episode)
	}

	return episodes, nil
}
