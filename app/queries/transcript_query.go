package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// getTranscript returns the parsed transcript for the user.
func getTranscript(scraper repository.ScraperProvider, collector *colly.Collector, params *models.TranscriptRequestBody) ([]models.Transcript, error) {
	// Create empty transcript
	var transcript []models.Transcript

	// Get initial page
	_, html, err := scraper.Navigate(collector, params.Base, repository.TRANSCRIPT_ROUTE)

	// Check for initial success
	if err != nil {
		return transcript, err
	}

	// Parse transcript HTML
	transcript = append(transcript, parsers.ParseTranscript(html))

	return transcript, nil
}
