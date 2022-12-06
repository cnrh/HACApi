package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// getSchedule returns the parsed schedule for the user.
func getSchedule(scraper repository.ScraperProvider, parser repository.ParserProvider, collector *colly.Collector, params models.ScheduleRequestBody) ([]models.Schedule, error) {
	// Create empty schedule
	var schedule []models.Schedule

	// Get initial page
	_, html, err := scraper.Navigate(collector, params.Base, repository.SCHEDULE_ROUTE)

	// Check for initial success
	if err != nil {
		return schedule, err
	}

	// Parse schedule HTML
	schedule = append(schedule, parser.ParseSchedule(html))

	return schedule, nil
}
