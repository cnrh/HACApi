package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

type Queries struct {
	Scraper repository.ScraperProvider
}

func (queries Queries) GetClasswork(collector *colly.Collector, params *models.ClassworkRequestBody) ([]models.Classwork, error) {
	return getClasswork(queries.Scraper, collector, params)
}

func (queries Queries) GetIPRAll(collector *colly.Collector, params *models.IprAllRequestBody) ([]models.IPR, error) {
	return getIPRAll(queries.Scraper, collector, params)
}

func (queries Queries) GetIPR(collector *colly.Collector, params *models.IprRequestBody) ([]models.IPR, error) {
	return getIPR(queries.Scraper, collector, params)
}

func (queries Queries) GetLogin(collector *colly.Collector, params *models.LoginRequestBody) ([]models.Login, error) {
	return getLogin(queries.Scraper, collector, params)
}

func (queries Queries) GetReportCard(collector *colly.Collector, params *models.ReportCardRequestBody) ([]models.ReportCard, error) {
	return getReportCard(queries.Scraper, collector, params)
}

func (queries Queries) GetSchedule(collector *colly.Collector, params *models.ScheduleRequestBody) ([]models.Schedule, error) {
	return getSchedule(queries.Scraper, collector, params)
}

func (queries Queries) GetTranscript(collector *colly.Collector, params *models.TranscriptRequestBody) ([]models.Transcript, error) {
	return getTranscript(queries.Scraper, collector, params)
}

func NewQueries(scraper repository.ScraperProvider) Queries {
	return Queries{Scraper: scraper}
}
