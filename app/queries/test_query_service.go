package queries

import (
	"errors"
	"math"

	"github.com/Threqt1/HACApi/app/models"
	"github.com/gocolly/colly"
)

// TestQuerier is the querier meant to be
// used while testing. It returns the default
// struct for every query.
type TestQuerier struct {
}

// Send back a slice of the same length as params.MarkingPeriods, or one if the length is 0.
func (queries TestQuerier) GetClasswork(collector *colly.Collector, params models.ClassworkRequestBody) ([]models.Classwork, error) {
	// Get length in the interval [1, 6].
	length := int(math.Max(1, float64(len(params.MarkingPeriods))))
	// Make the slice and return.
	return make([]models.Classwork, length), nil
}

func (queries TestQuerier) GetIPRAll(collector *colly.Collector, params models.IprAllRequestBody) ([]models.IPR, error) {
	return []models.IPR{{}}, nil
}

func (queries TestQuerier) GetIPR(collector *colly.Collector, params models.IprRequestBody) ([]models.IPR, error) {
	return []models.IPR{{}}, nil
}

func (queries TestQuerier) GetLogin(collector *colly.Collector, params models.LoginRequestBody) ([]models.Login, error) {
	return []models.Login{{}}, nil
}

func (queries TestQuerier) GetReportCard(collector *colly.Collector, params models.ReportCardRequestBody) ([]models.ReportCard, error) {
	return []models.ReportCard{{}}, nil
}

func (queries TestQuerier) GetSchedule(collector *colly.Collector, params models.ScheduleRequestBody) ([]models.Schedule, error) {
	return []models.Schedule{{}}, nil
}

func (queries TestQuerier) GetTranscript(collector *colly.Collector, params models.TranscriptRequestBody) ([]models.Transcript, error) {
	return []models.Transcript{{}}, nil
}

// NewTestQuerier makes a new test querier.
func NewTestQuerier() TestQuerier {
	return TestQuerier{}
}

// The error thrown by the TestErrrorQuerier.
var ErrorBadQuery = errors.New("bad query")

// TestErrorQuerier is a test querier which
// only returns errors for all of its
// queries.
type TestErrorQuerier struct {
}

func (queries TestErrorQuerier) GetClasswork(collector *colly.Collector, params models.ClassworkRequestBody) ([]models.Classwork, error) {
	return nil, ErrorBadQuery
}

func (queries TestErrorQuerier) GetIPRAll(collector *colly.Collector, params models.IprAllRequestBody) ([]models.IPR, error) {
	return nil, ErrorBadQuery
}

func (queries TestErrorQuerier) GetIPR(collector *colly.Collector, params models.IprRequestBody) ([]models.IPR, error) {
	return nil, ErrorBadQuery
}

func (queries TestErrorQuerier) GetLogin(collector *colly.Collector, params models.LoginRequestBody) ([]models.Login, error) {
	return nil, ErrorBadQuery
}

func (queries TestErrorQuerier) GetReportCard(collector *colly.Collector, params models.ReportCardRequestBody) ([]models.ReportCard, error) {
	return nil, ErrorBadQuery
}

func (queries TestErrorQuerier) GetSchedule(collector *colly.Collector, params models.ScheduleRequestBody) ([]models.Schedule, error) {
	return nil, ErrorBadQuery
}

func (queries TestErrorQuerier) GetTranscript(collector *colly.Collector, params models.TranscriptRequestBody) ([]models.Transcript, error) {
	return nil, ErrorBadQuery
}

// NewTestErrorQuerier makes a new test querier that
// always errors.
func NewTestErrorQuerier() TestErrorQuerier {
	return TestErrorQuerier{}
}
