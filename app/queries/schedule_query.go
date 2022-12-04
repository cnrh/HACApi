package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// GetSchedule accepts a collector and a base, and returns a parsed schedule.
func GetSchedule(server *repository.Server, collector *colly.Collector, base string) (models.Schedule, error) {
	// Create empty schedule
	var schedule models.Schedule

	// Get initial page
	_, html, err := server.Scraper.Navigate(collector, base, repository.SCHEDULE_ROUTE)

	// Check for initial success
	if err != nil {
		return schedule, err
	}

	// Parse schedule HTML
	schedule = parsers.ParseSchedule(html)

	return schedule, nil
}
