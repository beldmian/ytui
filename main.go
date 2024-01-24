package main

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"

	"log"

	"github.com/gocolly/colly"
	"github.com/ktr0731/go-fuzzyfinder"
)

const (
	channelURL = "https://yewtu.be/channel/"
	videoURL   = "https://www.youtube.com"
)

func main() {
	c := colly.NewCollector()

	f, err := os.Open("config.json")
	if err != nil {
		log.Panicln("Cannot open config file")
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		log.Panicln("Error in config file format")
	}
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		log.Panicln("Error unmarshaing config")
	}

	var channels []Channel
	c.OnHTML("body", func(e *colly.HTMLElement) {
		var channel Channel
		name := e.ChildText("div.channel-profile")
		channel.Name = name

		e.ForEach("div.pure-u-1.pure-u-md-1-4", func(i int, e *colly.HTMLElement) {
			var vid Video
			vid.URL = e.ChildAttr("a", "href")
			vid.Name = e.ChildText("a > p")
			channel.Videos = append(channel.Videos, vid)
		})

		channels = append(channels, channel)
	})
	for _, id := range config.Channels {
		c.Visit(channelURL + id)
	}
	ch, err := fuzzyfinder.Find(
		channels,
		func(i int) string {
			return channels[i].Name
		},
	)
	if err != nil {
		log.Panicln("Fuzzyfinder error")
	}
	videos := channels[ch].Videos
	vid, err := fuzzyfinder.Find(
		videos,
		func(i int) string {
			return videos[i].Name
		},
	)
	if err != nil {
		log.Panicln("Fuzzyfinder error")
	}
	cmd := exec.Command("mpv", videoURL+videos[vid].URL)
	cmd.Run()
}
