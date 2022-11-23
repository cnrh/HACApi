package queries

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
	"github.com/Threqt1/HACApi/app/queries/parsers"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/Threqt1/HACApi/pkg/utils"
	"github.com/gocolly/colly"
)

func GetClasswork(loginCollector *colly.Collector, base string, markingPeriods []int) ([]models.Classwork, error) {
	//Get initial page
	collector, html, err := utils.NavigateTo(loginCollector, base, repository.CLASSWORK_ROUTE)

	//Check for initial success
	if err != nil {
		return nil, err
	}

	//Determine the current marking period suffix
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

	//Get other necessary fields
	viewstate, _ := html.Find("input[name='__VIEWSTATE']").Attr("value")
	viewstategen, _ := html.Find("input[name='__VIEWSTATEGENERATOR']").Attr("value")
	eventvalidation, _ := html.Find("input[name='__EVENTVALIDATION']").Attr("value")

	//Make structs for initial getClassworkHTMLPages call
	formData := partialFormData{ViewState: viewstate, ViewStateGen: viewstategen, EventValidation: eventvalidation, Url: repository.CLASSWORK_ROUTE, Base: base}
	recievedInfo := recievedClassworkInfo{Html: html, Mp: currMarkingPer, Suffix: markingPerSuffix}

	//Generate classwork
	recievedClasswork, err := generateClasswork(collector, markingPeriods, recievedInfo, formData)

	if err != nil {
		return nil, err
	}

	return recievedClasswork, nil
}

// recievedClassworkInfo struct representing classwork information
// for a given marking period that was recieved by the first call.
type recievedClassworkInfo struct {
	Html   *goquery.Selection //The HTML for the recieved marking period information
	Mp     int                //The marking period the info is for
	Suffix string             //The suffix for MP selection
}

// getClassworkHTMLPages is a generator function which returns a channel where recieved HTML pages for each marking period's classwork will be sent.
func getClassworkHTMLPages(collector *colly.Collector, doneChan <-chan struct{}, markingPeriods []int, recievedInfo recievedClassworkInfo, formData partialFormData) chan PipelineResponse[*goquery.Selection] {
	//Output channels for recieved HTML
	htmlPagesChan := make(chan PipelineResponse[*goquery.Selection])

	//Wrap code in a goroutine so we can return channels
	go func() {
		var wg sync.WaitGroup

		//If no marking periods provided, return the current marking period data
		if len(markingPeriods) == 0 {
			select {
			case htmlPagesChan <- PipelineResponse[*goquery.Selection]{Value: recievedInfo.Html, Err: nil}:
			case <-doneChan:
			}
		}

		//Scrape in parallel for each goroutine
		for _, mp := range markingPeriods {
			wg.Add(1)

			go func(mp int) {
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

				if mp == recievedInfo.Mp {
					html = recievedInfo.Html
				} else {
					_, html, err = utils.PostTo(collector, formData.Base, formData.Url, utils.MakeClassworkFormData(strconv.Itoa(mp)+recievedInfo.Suffix, formData.ViewState, formData.ViewStateGen, formData.EventValidation))
				}

				//Send HTML/Error to output channel
				select {
				case htmlPagesChan <- PipelineResponse[*goquery.Selection]{Value: html, Err: err}:
				case <-doneChan:
				}
			}(mp)
		}

		//Wait till goroutines are done, then close channel
		go func() {
			wg.Wait()
			close(htmlPagesChan)
		}()
	}()

	return htmlPagesChan
}

// parseClassworkHTML parses the raw HTML outputted by the getClassworkHTML function, and returns a parsed Classwwork struct/propagates errors down.
func parseClassworkHTML(collector *colly.Collector, doneChan <-chan struct{}, markingPeriods []int, recievedInfo recievedClassworkInfo, formData partialFormData) chan PipelineResponse[models.Classwork] {
	//Make channel to emit parsed classwork and errors
	parsedClassworkChan := make(chan PipelineResponse[models.Classwork])

	go func() {
		//Recieve raw HTML
		rawHTMLChan := getClassworkHTMLPages(collector, doneChan, markingPeriods, recievedInfo, formData)

		var wg sync.WaitGroup

		//Parse each HTML recieved concurrently
		for htmlRes := range rawHTMLChan {
			//If error, cascade it down and break
			if htmlRes.Err != nil {
				parsedClassworkChan <- PipelineResponse[models.Classwork]{Value: models.Classwork{}, Err: htmlRes.Err}
				break
			}

			//Otherwise, start goroutine to parse
			wg.Add(1)
			go func(htmlRes PipelineResponse[*goquery.Selection]) {
				defer wg.Done()

				//Check if done was called beforehand
				select {
				case <-doneChan:
					return
				default:
				}

				//Parse classwork
				parsedClasswork := parsers.ParseClasswork(htmlRes.Value)

				//Try adding parsed classwork to channel
				select {
				case parsedClassworkChan <- PipelineResponse[models.Classwork]{Value: parsedClasswork, Err: nil}:
				case <-doneChan:
				}
			}(htmlRes)
		}

		//Wait until all parsing goroutines are done, then close outbound channel
		go func() {
			wg.Wait()
			close(parsedClassworkChan)
		}()
	}()

	return parsedClassworkChan
}

// generateClasswork oversees the raw HTML getter and raw HTML parser functions, and combines the outputted Classwork structs into an array, which is returned.
// Will also return an error, if there is one.
func generateClasswork(collector *colly.Collector, markingPeriods []int, recievedInfo recievedClassworkInfo, formData partialFormData) ([]models.Classwork, error) {
	//Make a done channel for cancellation
	doneChan := make(chan struct{})
	defer close(doneChan)

	//Recieve parsed classwork
	parsedClassworkChan := parseClassworkHTML(collector, doneChan, markingPeriods, recievedInfo, formData)

	//Append recieved classwork to array
	classworkArray := make([]models.Classwork, 0, len(markingPeriods))

	for parsedClasswork := range parsedClassworkChan {
		// Return if there was an error
		if parsedClasswork.Err != nil {
			return nil, parsedClasswork.Err
		}
		classworkArray = append(classworkArray, parsedClasswork.Value)
	}

	return classworkArray, nil
}
