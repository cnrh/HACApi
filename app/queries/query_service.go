package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

type Querier struct {
	Scraper repository.ScraperProvider
	Parser  repository.ParserProvider
}

func (queries Querier) GetClasswork(collector *colly.Collector, params models.ClassworkRequestBody) ([]models.Classwork, error) {
	return getClasswork(queries.Scraper, queries.Parser, collector, params)
}

func (queries Querier) GetIPRAll(collector *colly.Collector, params models.IprAllRequestBody) ([]models.IPR, error) {
	return getIPRAll(queries.Scraper, queries.Parser, collector, params)
}

func (queries Querier) GetIPR(collector *colly.Collector, params models.IprRequestBody) ([]models.IPR, error) {
	return getIPR(queries.Scraper, queries.Parser, collector, params)
}

func (queries Querier) GetLogin(collector *colly.Collector, params models.LoginRequestBody) ([]models.Login, error) {
	return getLogin(queries.Scraper, queries.Parser, collector, params)
}

func (queries Querier) GetReportCard(collector *colly.Collector, params models.ReportCardRequestBody) ([]models.ReportCard, error) {
	return getReportCard(queries.Scraper, queries.Parser, collector, params)
}

func (queries Querier) GetSchedule(collector *colly.Collector, params models.ScheduleRequestBody) ([]models.Schedule, error) {
	return getSchedule(queries.Scraper, queries.Parser, collector, params)
}

func (queries Querier) GetTranscript(collector *colly.Collector, params models.TranscriptRequestBody) ([]models.Transcript, error) {
	return getTranscript(queries.Scraper, queries.Parser, collector, params)
}

func NewQuerier(scraper repository.ScraperProvider, parser repository.ParserProvider) Querier {
	return Querier{Scraper: scraper, Parser: parser}
}
