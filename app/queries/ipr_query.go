package queries

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
)

// getIPR returns the latest IPR or the IPR for the date specified.
func getIPR(scraper repository.ScraperProvider, parser repository.ParserProvider, collector *colly.Collector, params models.IprRequestBody) ([]models.IPR, error) {
	// Get initial page
	collector, html, err := scraper.Navigate(collector, params.Base, repository.IPR_ROUTE)

	if err != nil {
		return nil, err
	}

	// Parse date
	date, err := time.Parse("01/02/2006", params.Date)

	if err != nil {
		return nil, err
	}

	// Determine current IPR date
	currDateOptionAttr := html.Find("#plnMain_ddlIPRDates > option[selected='selected']").Text()
	currDate, err := time.Parse("01/02/2006", currDateOptionAttr)
	if err != nil {
		return nil, err
	}

	// Get other necessary fields
	viewstate, _ := html.Find("input[name='__VIEWSTATE']").Attr("value")
	viewstategen, _ := html.Find("input[name='__VIEWSTATEGENERATOR']").Attr("value")
	eventvalidation, _ := html.Find("input[name='__EVENTVALIDATION']").Attr("value")

	// Make structs for pipeline generation
	formData := utils.PartialFormData{ViewState: viewstate, ViewStateGen: viewstategen, EventValidation: eventvalidation, Url: repository.IPR_ROUTE, Base: params.Base}
	recievedInfo := recievedIPRInfo{HTML: html, Date: currDate}
	functions := utils.PipelineFunctions[models.IPR, time.Time]{
		GenFormData: func(date string, pfd utils.PartialFormData) map[string]string {
			return utils.MakeIPRFormData(date, &pfd)
		},
		Parse: parser.ParseIPR,
		ToFormData: func(date time.Time) string {
			return date.Format("1/2/2006 03:04:05 PM")
		},
	}

	// Make array of dates
	dates := make([]time.Time, 0, 1)

	// If date isnt a zero value, append it into array
	if !date.IsZero() {
		dates = append(dates, date)
	}

	// Generate IPR
	recievedIPRs, err := utils.GeneratePipeline[models.IPR, time.Time](scraper, collector, dates, recievedInfo, &formData, functions)

	if err != nil {
		return nil, err
	}

	return recievedIPRs, nil
}

// recievedIPRInfo struct representing IPR information
// for a date that was recieved by the first call.
type recievedIPRInfo struct {
	HTML *goquery.Selection // The recieved HTML
	Date time.Time          // The time related to the IPR recieved
}

func (rii recievedIPRInfo) Html() *goquery.Selection {
	return rii.HTML
}

func (rii recievedIPRInfo) Equal(other time.Time) bool {
	return other.Equal(rii.Date)
}
