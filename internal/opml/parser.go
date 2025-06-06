package opml

import (
	"encoding/xml"
	"os"
)

// OPML represents the root structure of an OPML file
type OPML struct {
	XMLName xml.Name `xml:"opml"`
	Head    Head     `xml:"head"`
	Body    Body     `xml:"body"`
}

// Head represents the head section of an OPML file
type Head struct {
	Title string `xml:"title,omitempty"`
}

// Body represents the body section of an OPML file
type Body struct {
	Outlines []Outline `xml:"outline"`
}

// Outline represents an outline element in the OPML file
type Outline struct {
	Title  string `xml:"title,attr,omitempty"`
	XMLUrl string `xml:"xmlUrl,attr,omitempty"`
	Text   string `xml:"text,attr,omitempty"`
}

// ParseOPML parses an OPML file and returns the parsed structure
func ParseOPML(filename string) (*OPML, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var opml OPML
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&opml)
	if err != nil {
		return nil, err
	}

	return &opml, nil
}
