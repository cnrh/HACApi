package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

func GetTranscript(server *repository.Server, collector *colly.Collector, base string) (models.Transcript, error) {
	// Create empty transcript
	var transcript models.Transcript

	// Get initial page
	_, html, err := server.Scraper.Navigate(collector, base, repository.TRANSCRIPT_ROUTE)

	// Check for initial success
	if err != nil {
		return transcript, err
	}

	// Parse transcript HTML
	transcript = parsers.ParseTranscript(html)

	return transcript, nil
}
