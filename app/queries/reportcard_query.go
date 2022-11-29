package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
)

// GetReportCard accepts a collector and a base, and returns a parsed report card for the
// user logged into the collector.
func GetReportCard(loginCollector *colly.Collector, base string) (models.ReportCard, error) {
	// Create empty report card model
	var reportCard models.ReportCard

	// Get initial page
	_, html, err := utils.NavigateTo(loginCollector, base, repository.REPORT_CARD_ROUTE)

	// Check for initial success
	if err != nil {
		return reportCard, err
	}

	// Parse report card HTML
	reportCard = parsers.ParseReportCard(html)

	return reportCard, nil
}
