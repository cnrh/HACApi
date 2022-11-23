package queries

import (
	"errors"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
)

func GetIPR(loginCollector *colly.Collector, base string, date time.Time) ([]models.IPR, error) {
	//Get initial page
	collector, html, err := utils.NavigateTo(loginCollector, base, repository.IPR_ROUTE)

	//Check for initial success
	if err != nil {
		return nil, err
	}

	//Determine current IPR date

	currDateOptionAttr, exists := html.Find("#plnMain_ddlIPRDates > option[selected='selected']").Attr("value")
	if !exists {
		return nil, errors.New("invalid page")
	}
	currDate, err := time.Parse("1/2/2006 03:04:05 PM", currDateOptionAttr)
	if err != nil {
		return nil, err
	}

	//Get other necessary fields
	viewstate, _ := html.Find("input[name='__VIEWSTATE']").Attr("value")
	viewstategen, _ := html.Find("input[name='__VIEWSTATEGENERATOR']").Attr("value")
	eventvalidation, _ := html.Find("input[name='__EVENTVALIDATION']").Attr("value")

	//Make structs for initial getIPRHTMLPages call
	formData := partialFormData{ViewState: viewstate, ViewStateGen: viewstategen, EventValidation: eventvalidation, Url: repository.IPR_ROUTE, Base: base}
	recievedInfo := recievedIPRInfo{html: html, date: currDate}

	//Generate IPR
	recievedIPRs, err := generateIPRs(collector, []time.Time{date}, recievedInfo, formData)

	if err != nil {
		return nil, err
	}

	return recievedIPRs, nil
}

// recievedIPRInfo struct representing IPR information
// for a date that was recieved by the first call.
type recievedIPRInfo struct {
	html *goquery.Selection //The recieved HTML
	date time.Time          //The time related to the IPR recieved
}

// getIPRHTMLPages is a generator function which returns a channel where raw HTML pages of IPRs for each date inputted will be sent
func getIPRHTMLPages(collector *colly.Collector, doneChan <-chan struct{}, dates []time.Time, recievedInfo recievedIPRInfo, formData partialFormData) chan PipelineResponse[*goquery.Selection] {
	//Output channels for recieved HTML
	htmlPagesChan := make(chan PipelineResponse[*goquery.Selection])

	//Wrap code in a goroutine so we can return channels
	go func() {
		var wg sync.WaitGroup

		//If no dates are inputted, return most recent IPR
		if len(dates) == 0 {
			select {
			case htmlPagesChan <- PipelineResponse[*goquery.Selection]{Value: recievedInfo.html, Err: nil}:
			case <-doneChan:
			}
		}

		//Scrape in parallel for each date
		for _, date := range dates {
			wg.Add(1)

			go func(date time.Time) {
				defer wg.Done()

				//Check if done's been called before expensive request
				select {
				case <-doneChan:
					return
				default:
				}

				//Get HTML
				var html *goquery.Selection
				var err error

				if date.Equal(recievedInfo.date) {
					html = recievedInfo.html
				} else {
					//Format date to appropriate format
					formattedDate := date.Format("1/2/2006 03:04:05 PM")
					_, html, err = utils.PostTo(collector, formData.Base, formData.Url, utils.MakeIPRFormData(formattedDate, formData.ViewState, formData.ViewStateGen, formData.EventValidation))
				}

				//Send HTML/Error to output channel
				select {
				case htmlPagesChan <- PipelineResponse[*goquery.Selection]{Value: html, Err: err}:
				case <-doneChan:
				}
			}(date)
		}

		//Wait till goroutines are done, then close channel
		go func() {
			wg.Wait()
			close(htmlPagesChan)
		}()
	}()

	return htmlPagesChan
}

func parseIPRHTML(collector *colly.Collector, doneChan <-chan struct{}, dates []time.Time, recievedInfo recievedIPRInfo, formData partialFormData) chan PipelineResponse[models.IPR] {
	//Make a channel to output IPR
	parsedIPRChan := make(chan PipelineResponse[models.IPR])

	go func() {
		//Recieve raw html
		rawHTMLChan := getIPRHTMLPages(collector, doneChan, dates, recievedInfo, formData)

		var wg sync.WaitGroup

		//Parse HTML concurrently
		for htmlRes := range rawHTMLChan {
			//If error, cascade down and break
			if htmlRes.Err != nil {
				parsedIPRChan <- PipelineResponse[models.IPR]{Value: models.IPR{}, Err: htmlRes.Err}
				break
			}

			//Otherwise, start parsing
			wg.Add(1)
			go func(htmlRes PipelineResponse[*goquery.Selection]) {
				defer wg.Done()

				//Check if done was called beforehand
				select {
				case <-doneChan:
					return
				default:
				}

				//Parse IPR
				parsedIPR := parsers.ParseIPR(htmlRes.Value)

				//Try adding parsed IPR to channel
				select {
				case parsedIPRChan <- PipelineResponse[models.IPR]{Value: parsedIPR, Err: nil}:
				case <-doneChan:
				}
			}(htmlRes)
		}

		//Clean up channels after work's done
		go func() {
			wg.Wait()
			close(parsedIPRChan)
		}()
	}()

	return parsedIPRChan
}

func generateIPRs(collector *colly.Collector, dates []time.Time, recievedInfo recievedIPRInfo, formData partialFormData) ([]models.IPR, error) {
	//Make done channel for cancellation
	doneChan := make(chan struct{})
	defer close(doneChan)

	//Recieved parsed IPRs
	parsedIPRChan := parseIPRHTML(collector, doneChan, dates, recievedInfo, formData)

	//Append recieved IPRs to array
	iprArray := make([]models.IPR, 0, len(dates))

	for parsedIPR := range parsedIPRChan {
		//Return error if error
		if parsedIPR.Err != nil {
			return nil, parsedIPR.Err
		}
		iprArray = append(iprArray, parsedIPR.Value)
	}

	return iprArray, nil
}
