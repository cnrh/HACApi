package queries

import (
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
)

func GetReportCard(loginCollector *colly.Collector, base string) (models.ReportCard, error) {
	//Create report card model
	reportCard := models.ReportCard{}

	//Get initial page
	_, html, err := utils.NavigateTo(loginCollector, base, repository.REPORT_CARD_ROUTE)

	//Check for initial success
	if err != nil {
		return reportCard, err
	}

	//Parse report card HTML
	recievedReportCard := parsers.ParseReportCard(html)

	return recievedReportCard, nil
}
