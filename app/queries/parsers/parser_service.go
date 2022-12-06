package parsers

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
)

type Parser struct {
}

func (parser Parser) ParseClasswork(html *goquery.Selection) models.Classwork {
	return parseClasswork(html)
}

func (parser Parser) ParseIPR(html *goquery.Selection) models.IPR {
	return parseIPR(html)
}

func (parser Parser) ParseReportCard(html *goquery.Selection) models.ReportCard {
	return parseReportCard(html)
}

func (parser Parser) ParseSchedule(html *goquery.Selection) models.Schedule {
	return parseSchedule(html)
}

func (parser Parser) ParseTranscript(html *goquery.Selection) models.Transcript {
	return parseTranscript(html)
}

func NewParser() Parser {
	return Parser{}
}
