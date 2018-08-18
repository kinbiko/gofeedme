package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/bugsnag/bugsnag-go"
	"github.com/mmcdole/gofeed"
)

const (
	// bugsnagAPIKey is the API key for this project on Bugsnag.
	bugsnagAPIKey = "47f6b13fc2258ad2d4b1a78766fe00ba"

	// version is the current gofeedme
	// TODO: extract this from the version of the Github repo when compiling
	version = "1.0.0"

	// widgetFilename is the filename of the config file when used as an Übersicht plugin
	widgetFilename = "GoFeedMe.widget/config.json"

	// localFilename is the filename of the config file when testing locally
	localFilename = "./config.json"

	// itemCount is the number of RSS feed items to fetch for each source
	itemCount = 3
)

var filename = widgetFilename

type feed struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type configuration struct {
	Feeds []feed `json:"feeds"`
}

func main() {
	configureBugsnag()
	defer bugsnag.AutoNotify()
	if _, err := http.Get("https://google.com/"); err != nil {
		//Chances are if Google's down, the internet is down.
		fmt.Println("No network connection. Cannot get RSS feed")
		return
	}

	for _, configFeed := range parseConfig(readConfigFile()).Feeds {
		fetchFeed(configFeed)
	}
}

func configureBugsnag() {
	bugsnag.Configure(bugsnag.Configuration{
		APIKey:              bugsnagAPIKey,
		AppType:             "Übersicht",
		AppVersion:          version,
		NotifyReleaseStages: []string{"production"},
		ProjectPackages:     []string{"main", "github.com/kinbiko/*"},
		ReleaseStage:        "development",                   //TODO: turn into a build flag
		SourceRoot:          "/Users/kinbiko/repos/gofeedme", //TODO extract programatically
	})
}

func readConfigFile() []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return bytes
}

func parseConfig(bytes []byte) configuration {
	conf := configuration{}
	err := json.Unmarshal(bytes, &conf)
	if err != nil {
		panic(fmt.Sprintf("unable to unmarshal file to configuration: %v", err))
	}
	return conf
}

func fetchFeed(configFeed feed) {
	parsedFeed, err := gofeed.NewParser().ParseURL(configFeed.URL)
	if err != nil {
		bugsnag.Notify(fmt.Errorf("unable to parse feed, ignoring '%s': %v", configFeed.Name, err))
		return
	}

	fmt.Println("<h4>" + configFeed.Name + "</h4>")
	fmt.Println(makeLinks(parsedFeed.Items))
}

func makeLinks(stories []*gofeed.Item) string {
	str := "<ul>"
	for _, story := range stories[:itemCount] {
		str += "<li><a href='" + story.Link + "'>" + story.Title + "</a></li>"
	}
	return str + "</ul>"
}
