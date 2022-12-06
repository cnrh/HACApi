package queries

import (
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
)

// getIPRAll returns all the IPRs registered for the user, or the dates only if specified.
func getIPRAll(scraper repository.ScraperProvider, collector *colly.Collector, params *models.IprAllRequestBody) ([]models.IPR, error) {
	// Get initial page
	collector, html, err := scraper.Navigate(collector, params.Base, repository.IPR_ROUTE)

	// Check for initial success
	if err != nil {
		return nil, err
	}

	// Determine current IPR date
	currDateOptionAttr := html.Find("#plnMain_ddlIPRDates > option[selected='selected']").Text()
	currDate, err := time.Parse("01/02/2006", currDateOptionAttr)
	if err != nil {
		return nil, err
	}

	// Get every single avaliable date
	dateOptionEles := html.Find("#plnMain_ddlIPRDates > option")
	dates := make([]time.Time, 0, dateOptionEles.Length())

	dateOptionEles.Each(func(_ int, dateOptionEle *goquery.Selection) {
		// Get text
		dateText := dateOptionEle.Text()

		// Parse date
		date, err := time.Parse("01/02/2006", dateText)

		// If no err, push to dates
		if err == nil {
			dates = append(dates, date)
		}
	})

	// If only dates were needed, convert dates into correct model and return
	if params.DatesOnly {
		partialIPRs := make([]models.IPR, 0, len(dates))
		for _, date := range dates {
			partialIPRs = append(partialIPRs, models.IPR{Date: date.Format("01/02/2006"), Entries: []models.IPREntry{}})
		}
		return partialIPRs, nil
	}

	// Get other necessary fields
	viewstate, _ := html.Find("input[name='__VIEWSTATE']").Attr("value")
	viewstategen, _ := html.Find("input[name='__VIEWSTATEGENERATOR']").Attr("value")
	eventvalidation, _ := html.Find("input[name='__EVENTVALIDATION']").Attr("value")

	// Make structs for pipeline generation
	formData := utils.PartialFormData{ViewState: viewstate, ViewStateGen: viewstategen, EventValidation: eventvalidation, Url: repository.IPR_ROUTE, Base: params.Base}
	recievedInfo := recievedIPRInfo{HTML: html, Date: currDate}
	functions := utils.PipelineFunctions[models.IPR, time.Time]{
		GenFormData: utils.MakeIPRFormData,
		Parse:       parsers.ParseIPR,
		ToFormData: func(date time.Time) string {
			return date.Format("1/2/2006 03:04:05 PM")
		},
	}

	// Generate IPRs
	recievedIPRs, err := utils.GeneratePipeline[models.IPR, time.Time](scraper, collector, dates, recievedInfo, formData, functions)

	if err != nil {
		return nil, err
	}

	return recievedIPRs, nil
}
