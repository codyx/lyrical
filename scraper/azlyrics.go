package scraper

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/codyx/lyrical/prompt"
	"github.com/gocolly/colly"
)

const (
	// AZBaseURL is the URL used to retrieve the lyrics of a song.
	AZBaseURL = "https://www.azlyrics.com"

	// AZSearchURL is the URL used to search for a song.
	AZSearchURL = "https://search.azlyrics.com"
)

// AZLyrics is the scraper struct for AZLyrics.
type AZLyrics struct {
	Collector *colly.Collector
	Prompt    *prompt.Prompt
	Lyrics    string
}

// New instantiates a new AZLyrics and returns a pointer to it.
func New() *AZLyrics {
	az := &AZLyrics{
		Collector: colly.NewCollector(),
		Prompt:    prompt.New(),
	}
	return az
}

func urlEncode(s string, d string) string {
	return strings.ToLower(strings.Replace(url.QueryEscape(s), "+", d, -1))
}

// SearchSongURL is the advanced search for a song in case there is no match found
// by SearchSongLyrics.
func (a *AZLyrics) SearchSongURL() {
	url := fmt.Sprintf("%s/search.php?q=%s", AZSearchURL,
		urlEncode(a.Prompt.ArtistName+" "+a.Prompt.SongName, "+"))

	a.Collector.OnHTML("td > a[href]", func(e *colly.HTMLElement) {
		name, link := e.Text, e.Attr("href")
		if strings.HasPrefix(link, "http") {
			if _, ok := a.Prompt.SearchResult[name]; !ok {
				a.Prompt.SearchResult[name] = link
			}
		}
	})

	a.Collector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Search failed with status code:", r.StatusCode)
	})

	a.Collector.Visit(url)
}

// SearchSongLyrics starts a new lyrics search with a prompt.
func (a *AZLyrics) SearchSongLyrics() {
	a.Prompt.Run()

	artist, song := a.Prompt.ArtistName, a.Prompt.SongName
	url := fmt.Sprintf("%s/lyrics/%s/%s.html", AZBaseURL, urlEncode(artist, ""), urlEncode(song, ""))

	a.Collector.OnHTML("div:not([class])", func(e *colly.HTMLElement) {
		a.Lyrics = strings.TrimLeft(e.Text, "\n")
		fmt.Printf("%s", a.Lyrics)
	})

	a.Collector.OnScraped(func(r *colly.Response) {
		resultsNb := len(a.Prompt.SearchResult)
		if a.Lyrics == "" && resultsNb != 0 {
			link := a.Prompt.SelectSearchResult()
			if link != "" {
				a.Collector.OnError(func(r *colly.Response, err error) {
					if r.StatusCode == 404 {
						a.SearchSongURL()
					}
				})

				a.Collector.Visit(link)
			}
		} else if a.Lyrics == "" && resultsNb == 0 {
			fmt.Println("No song found")
		}
	})

	a.Collector.OnError(func(r *colly.Response, err error) {
		if r.StatusCode == 404 {
			a.SearchSongURL()
		}
	})

	a.Collector.Visit(url)
}
