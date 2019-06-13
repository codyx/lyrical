package prompt

import (
	"os"

	"github.com/manifoldco/promptui"
)

// Prompt is the entity representing
// the current search metadata.
type Prompt struct {
	ArtistName   string
	SongName     string
	SearchResult map[string]string
}

// New instantiates a new prompt and returns a pointer to it.
func New() *Prompt {
	return &Prompt{SearchResult: make(map[string]string)}
}

func checkErr(err error) {
	if err != nil {
		os.Exit(1)
	}
}

// Run the prompt and save the artist and song name.
func (p *Prompt) Run() {
	labels := []string{"Artist", "Song"}
	for _, label := range labels {
		prompt := promptui.Prompt{
			Label: label,
		}

		result, err := prompt.Run()

		checkErr(err)

		switch label {
		case "Artist":
			p.ArtistName = result
		case "Song":
			p.SongName = result
		}
	}
}

// SelectSearchResult creates a select prompt to help
// the user choose a song among a list of results.
func (p *Prompt) SelectSearchResult() string {
	results := []string{}

	for songName := range p.SearchResult {
		results = append(results, songName)
	}

	prompt := promptui.Select{
		Label: "Select a Result",
		Items: results,
	}

	_, result, err := prompt.Run()
	checkErr(err)

	return p.SearchResult[result]
}
