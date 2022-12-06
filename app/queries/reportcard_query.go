package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// getReportCard returns the parsed report card for the user.
func getReportCard(scraper repository.ScraperProvider, collector *colly.Collector, params *models.ReportCardRequestBody) ([]models.ReportCard, error) {
	// Create empty report card model
	var reportCard []models.ReportCard

	// Get initial page
	_, html, err := scraper.Navigate(collector, params.Base, repository.REPORT_CARD_ROUTE)

	// Check for initial success
	if err != nil {
		return reportCard, err
	}

	// Parse report card HTML
	reportCard = append(reportCard, parsers.ParseReportCard(html))

	return reportCard, nil
}
