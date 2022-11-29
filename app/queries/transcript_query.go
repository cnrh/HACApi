package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
)

func GetTranscript(loginCollector *colly.Collector, base string) (models.Transcript, error) {
	// Create empty transcript
	var transcript models.Transcript

	// Get initial page
	_, html, err := utils.NavigateTo(loginCollector, base, repository.TRANSCRIPT_ROUTE)

	// Check for initial success
	if err != nil {
		return transcript, err
	}

	// Parse transcript HTML
	transcript = parsers.ParseTranscript(html)

	return transcript, nil
}
