package main

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
)

const FILENAME = "GoFeedMe.widget/config.json"
const ITEMS_COUNT = 3

type feed struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type config struct {
	Feeds []feed `json:"feeds"`
}

func main() {
	config := parseConfig()
	for _, configFeed := range config.Feeds {
		parsedFeed, err := gofeed.NewParser().ParseURL(configFeed.URL)
		if err != nil {
			fmt.Println("unable to parse feed, ignoring"+configFeed.Name, err)
			continue
		}

		printHeader(configFeed, parsedFeed)
		printLinks(parsedFeed)
	}
}

func printHeader(configFeed feed, rssFeed *gofeed.Feed) {
	var name string
	if configFeed.Name != "" {
		name = configFeed.Name
	} else {
		name = rssFeed.Title
	}

	fmt.Printf("<h4>%s</h4>\n", name)
}

func printLinks(rssFeed *gofeed.Feed) {
	fmt.Println("<ul>")
	for _, story := range rssFeed.Items[:ITEMS_COUNT] {
		fmt.Printf("<li>")
		fmt.Printf("<a href='%s'>", story.Link)
		fmt.Printf(story.Title)
		fmt.Printf("</a>")
		fmt.Println("</li>")
	}
	fmt.Println("</ul>")
}

func parseConfig() config {
	bytes, err := ioutil.ReadFile(FILENAME)
	if err != nil {
		fmt.Println("unable to parse file:", err)
		panic(err)
	}
	conf := config{}
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		fmt.Println("unable to unmarshal file to config", err)
		panic(err)
	}
	return conf
}
