package main

import (
	"github.com/codyx/lyrical/scraper"
)

func main() {
	s := scraper.New()

	s.SearchSongLyrics()
}
