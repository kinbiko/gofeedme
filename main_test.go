package main

import (
	"testing"

	bugsnag "github.com/bugsnag/bugsnag-go"
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

func TestConfiguresBugsnagCorrectly(t *testing.T) {
	configureBugsnag()
	c := bugsnag.Config
	tt := []struct {
		property string
		exp      interface{}
		got      interface{}
	}{
		{property: "APIKey", exp: bugsnagAPIKey, got: c.APIKey},
		{property: "AppType", exp: "Übersicht", got: c.AppType},
		{property: "AppVersion", exp: version, got: c.AppVersion},
		{property: "ReleaseStage", exp: "development", got: c.ReleaseStage},
		{property: "SourceRoot", exp: "/Users/kinbiko/repos/gofeedme", got: c.SourceRoot},
	}
	for _, tc := range tt {
		if tc.got != tc.exp {
			t.Errorf("Expected bugsnag config parameter '%s' to be '%s' but was '%s'", tc.property, tc.exp, tc.got)
		}
	}
}
