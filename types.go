package main

type Channel struct {
	Name   string
	Videos []Video
}

type Video struct {
	URL  string
	Name string
}

type Config struct {
	Channels []string
}
