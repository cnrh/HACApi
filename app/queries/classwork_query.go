package queries

import (
	"errors"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
)

// GetClasswork accepts a collector, base, and a list of marking periods, and returns an array of parsed classwork
// objects for the inputted marking periods.
func GetClasswork(server *repository.Server, collector *colly.Collector, base string, markingPeriods []int) ([]models.Classwork, error) {
	// Get initial page
	collector, html, err := server.Scraper.Navigate(collector, base, repository.CLASSWORK_ROUTE)

	// Check for initial success
	if err != nil {
		return nil, err
	}

	// Determine the current marking period suffix
	markingPerOptionAttr, exists := html.Find("#plnMain_ddlReportCardRuns > option[selected='selected']").Attr("value")
	if !exists {
		return nil, errors.New("invalid page")
	}
	markingPerOptionText := strings.TrimSpace(markingPerOptionAttr)
	markingPerSuffix := markingPerOptionText[1:]
	currMarkingPer, err := strconv.Atoi(markingPerOptionText[0:1])
	if err != nil {
		return nil, err
	}

	// Get other necessary fields
	viewstate, _ := html.Find("input[name='__VIEWSTATE']").Attr("value")
	viewstategen, _ := html.Find("input[name='__VIEWSTATEGENERATOR']").Attr("value")
	eventvalidation, _ := html.Find("input[name='__EVENTVALIDATION']").Attr("value")

	// Make structs for pipeline generation
	formData := utils.PartialFormData{ViewState: viewstate, ViewStateGen: viewstategen, EventValidation: eventvalidation, Url: repository.CLASSWORK_ROUTE, Base: base}
	recievedInfo := recievedClassworkInfo{HTML: html, Mp: currMarkingPer}
	functions := utils.PipelineFunctions[models.Classwork, int]{
		GenFormData: utils.MakeClassworkFormData,
		Parse:       parsers.ParseClasswork,
		ToFormData: func(mp int) string {
			return strconv.Itoa(mp) + markingPerSuffix
		},
	}

	if len(markingPeriods) == 0 {
		markingPeriods = append(markingPeriods, recievedInfo.Mp)
	}

	// Generate classwork
	recievedClasswork, err := utils.GeneratePipeline[models.Classwork, int](server, collector, markingPeriods, recievedInfo, formData, functions)

	if err != nil {
		return nil, err
	}

	return recievedClasswork, nil
}

// recievedClassworkInfo struct representing classwork information
// for a given marking period that was recieved by the first call.
type recievedClassworkInfo struct {
	HTML *goquery.Selection // The recieved HTML
	Mp   int                // The marking period
}

func (rci recievedClassworkInfo) Html() *goquery.Selection {
	return rci.HTML
}

func (rci recievedClassworkInfo) Equal(other int) bool {
	return rci.Mp == other
}
