package queries

import (
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

func GetClasswork(loginCollector *colly.Collector, base string, markingPeriods []int) []models.Classwork {
	//Get initial page
	collector, html, err := utils.NavigateTo(loginCollector, base, repository.CLASSWORK_ROUTE)

	//Check for initial success
	if err != nil {
		return nil
	}

	//Determine the current marking period suffix
	var currMarkingPer int

	markingPerOptionAttr, exists := html.Find("#plnMain_ddlReportCardRuns > option[selected='selected']").Attr("value")
	if !exists {
		return nil
	}
	markingPerOptionText := strings.TrimSpace(markingPerOptionAttr)
	markingPerSuffix := markingPerOptionText[1:]
	currMarkingPerConv, err := strconv.Atoi(markingPerOptionText[0:1])
	if err == nil {
		currMarkingPer = currMarkingPerConv
	}

	//Get other necessary fields
	viewstate, _ := html.Find("input[name='__VIEWSTATE']").Attr("value")
	viewstategen, _ := html.Find("input[name='__VIEWSTATEGENERATOR']").Attr("value")
	eventvalidation, _ := html.Find("input[name='__EVENTVALIDATION']").Attr("value")

	//Make structs for initial getHTMLPages call
	partialFormData := partialClassworkFormData{ViewState: viewstate, ViewStateGen: viewstategen, EventValidation: eventvalidation, Url: repository.CLASSWORK_ROUTE, Base: base}
	recievedSwInfo := recievedClassworkMpInfo{Html: html, Mp: currMarkingPer, Suffix: markingPerSuffix}

	//Generate classwork
	recievedClasswork, err := generateClasswork(collector, markingPeriods, recievedSwInfo, partialFormData)

	if err != nil {
		return nil
	}

	return recievedClasswork
}

// partialClassworkFormData struct representing partial formdata
// for the classwork POST request.
type partialClassworkFormData struct {
	ViewState       string //viewstate formdata entry
	ViewStateGen    string //viewstategen formdata entry
	EventValidation string //eventvalidation formdata entry
	Url             string //url for the request
	Base            string //base url for the request
}

// recievedClassworkMpInfo struct representing classwork information
// for a given marking period that was recieved by the first call.
type recievedClassworkMpInfo struct {
	Html   *goquery.Selection //The HTML for the recieved marking period information
	Mp     int                //The marking period the info is for
	Suffix string             //The suffix for MP selection
}

// getClassworkHTMLPages is a generator function which returns a channel where recieved HTML pages for each marking period's classwork will be sent.
func getClassworkHTMLPages(collector *colly.Collector, doneChan <-chan struct{}, markingPeriods []int, recievedMpInfo recievedClassworkMpInfo, partialFormData partialClassworkFormData) chan PipelineResponse[*goquery.Selection] {
	//Output channels for recieved HTML
	htmlPagesChan := make(chan PipelineResponse[*goquery.Selection])

	//Wrap code in a goroutine so we can return channels
	go func() {
		var wg sync.WaitGroup

		//If no marking periods provided, return the current marking period data
		if len(markingPeriods) == 0 {
			select {
			case htmlPagesChan <- PipelineResponse[*goquery.Selection]{Value: recievedMpInfo.Html, Err: nil}:
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
				if mp == recievedMpInfo.Mp {
					html = recievedMpInfo.Html
				} else {
					_, html, err = utils.PostTo(collector, partialFormData.Base, partialFormData.Url, utils.MakeClassworkFormData(strconv.Itoa(mp)+recievedMpInfo.Suffix, partialFormData.ViewState, partialFormData.ViewStateGen, partialFormData.EventValidation))
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

func parseRawClassworkHTML(collector *colly.Collector, doneChan <-chan struct{}, markingPeriods []int, recievedMpInfo recievedClassworkMpInfo, partialFormData partialClassworkFormData) chan PipelineResponse[models.Classwork] {
	//Make channel to emit parsed classwork and errors
	parsedClassworkChan := make(chan PipelineResponse[models.Classwork])

	go func() {
		//Recieve raw HTML
		rawHTMLChan := getClassworkHTMLPages(collector, doneChan, markingPeriods, recievedMpInfo, partialFormData)

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
func generateClasswork(collector *colly.Collector, markingPeriods []int, recievedMpInfo recievedClassworkMpInfo, partialFormData partialClassworkFormData) ([]models.Classwork, error) {
	//Make a done channel for cancellation
	doneChan := make(chan struct{})
	defer close(doneChan)

	//Recieve parsed HTML pages
	parsedClassworkChan := parseRawClassworkHTML(collector, doneChan, markingPeriods, recievedMpInfo, partialFormData)

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
