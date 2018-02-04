package main

import (
	"encoding/json"
	"fmt"
	"github.com/mmcdole/gofeed"
	"io/ioutil"
	"net/http"
)

// WIDGET_FILENAME is the filename of the config file when used as an ubersicht plugin
const WIDGET_FILENAME = "GoFeedMe.widget/config.json"

// LOCAL_FILENAME is the filename of the config file when testing locally
const LOCAL_FILENAME = "./config.json"

// ITEMS_COUNT is the number of RSS feed items to fetch for each source
const ITEMS_COUNT = 3

const filename = WIDGET_FILENAME

type feed struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type config struct {
	Feeds []feed `json:"feeds"`
}

func main() {
	if offline() {
		fmt.Println("No network connection. Cannot get RSS feed")
		return
	}
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

func offline() bool {
	//Chances are if Google's down, the internet is down.
	_, err := http.Get("http://google.com/")
	return err != nil
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
	bytes, err := ioutil.ReadFile(filename)
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
