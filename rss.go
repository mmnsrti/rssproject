package main

import (
	"encoding/xml"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Link        string    `xml:"link"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}
type RSSItem struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (*RSSFeed, error) {

	httpClient := &http.Client{
		Timeout: 10 * 1000, // 10 seconds
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat ,err := io.ReadAll(resp.Body) // Read the response body to ensure the request is complete
	if err != nil {
		return nil, err
	}
	rssFeed:= RSSFeed{}
	xml.Unmarshal(dat, &rssFeed)
}
