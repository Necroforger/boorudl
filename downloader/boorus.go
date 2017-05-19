package downloader

import "net/http"

// Only store booru structs here.
// This file will most likely be used with go:generate in the future

// Danbooru searches images from http://danbooru.donmai.us/
type Danbooru struct {
	client *http.Client
}

// Gelbooru searches images from https://gelbooru.com/
type Gelbooru struct {
	client *http.Client
}

// Rule34 searches images from http://rule34.xxx/
type Rule34 struct {
	client *http.Client
}
