package main

import (
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestParsesConfig(t *testing.T) {
	config := parseConfig([]byte(`{
	"feeds": [
		{
			"name": "Kinbiko Blog",
			"url": "https://kinbiko.com/rss.xml"
		}
	]
}`))
	if expected, got := 1, len(config.Feeds); expected != got {
		t.Errorf("Expected %d feed(s) but config contains %d feed(s)", expected, got)
	}
	if expected, got := "Kinbiko Blog", config.Feeds[0].Name; expected != got {
		t.Errorf("Expected Feed name '%s' but was '%s'", expected, got)
	}
	if expected, got := "https://kinbiko.com/rss.xml", config.Feeds[0].URL; expected != got {
		t.Errorf("Expected Feed URL '%s' but was '%s'", expected, got)
	}
}

func TestMakeLinks(t *testing.T) {
	expected := "<ul><li><a href='a.com'>A</a></li><li><a href='b.com'>B</a></li><li><a href='c.com'>C</a></li></ul>"
	feeds := []*gofeed.Item{
		{Title: "A", Link: "a.com"},
		{Title: "B", Link: "b.com"},
		{Title: "C", Link: "c.com"},
	}
	got := makeLinks(feeds)
	if got != expected {
		t.Errorf("Expected links to be '%s' but were '%s'", expected, got)
	}
}
