package utils

import (
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// PipelineResponse represents a intermediary
// response between pipelines, used extensively
// in the queries package.
type PipelineResponse[T any] struct {
	Value T
	Err   error
}

// PartialFormData represents partial formdata
// for a POST request.
type PartialFormData struct {
	ViewState       string //viewstate formdata entry
	ViewStateGen    string //viewstategen formdata entry
	EventValidation string //eventvalidation formdata entry
	Url             string //url for the request
	Base            string //base url for the request
}

// PipelineFunctions describes the functions needed in order
// for the pipeline to correctly utilize the provided information.
type PipelineFunctions[T any, V any] struct {
	GenFormData func(string, PartialFormData) map[string]string
	Parse       func(*goquery.Selection) T
	ToFormData  func(V) string
}

// IPipleineRecievedValue represents an interface which
// contains the needed methods for the pipeline to correctly
// handle the recievedInformation parameter.
type IPipelineRecievedValue[V any] interface {
	Html() *goquery.Selection
	Equal(V) bool
}

// GeneratePipeline creates a new pipeline which will gather data from a POST request, parse it, and return it in an array format. T represents the model struct,
// V represents the recieved value's type.
func GeneratePipeline[T any, V any](collector *colly.Collector, data []V, recievedInfo IPipelineRecievedValue[V], formData PartialFormData, functions PipelineFunctions[T, V]) ([]T, error) {
	//Make a done channel for cancelling on error
	doneChan := make(chan struct{})
	defer close(doneChan)

	//Recieve parsed data
	parsedDataChan := pipelineParseHTML(collector, doneChan, data, recievedInfo, formData, functions)

	//Append recieved data to array
	dataArray := make([]T, 0, len(data))

	for res := range parsedDataChan {
		//If there's an error, cancel
		if res.Err != nil {
			return nil, res.Err
		}

		dataArray = append(dataArray, res.Value)
	}

	return dataArray, nil
}

// pipelineParseHTML represents the step in the pipeline where raw HTML is recieved through a channel, parsed, and emitted out through another channel.
func pipelineParseHTML[T any, V any](collector *colly.Collector, doneChan <-chan struct{}, data []V, recievedInfo IPipelineRecievedValue[V], formData PartialFormData, functions PipelineFunctions[T, V]) chan PipelineResponse[T] {
	//Make a channel to emit parsed data/errors
	parsedDataChan := make(chan PipelineResponse[T])

	go func() {
		//Recieve raw HTML
		rawHTMLChan := pipelineGetHTML(collector, doneChan, data, recievedInfo, formData, functions)

		var wg sync.WaitGroup

		//Parse HTML concurrently
		for res := range rawHTMLChan {
			//If error, cascade it down and break
			if res.Err != nil {
				parsedDataChan <- PipelineResponse[T]{Err: res.Err}
				break
			}

			//Otherwise, start goroutine to parse
			wg.Add(1)
			go func(res PipelineResponse[*goquery.Selection]) {
				defer wg.Done()

				//Check if done was called
				select {
				case <-doneChan:
					return
				default:
				}

				//Parse HTML
				parsedData := functions.Parse(res.Value)

				//Try emitting parsed data
				select {
				case parsedDataChan <- PipelineResponse[T]{Value: parsedData, Err: nil}:
				case <-doneChan:
				}
			}(res)
		}

		//Close channel once finished
		go func() {
			wg.Wait()
			close(parsedDataChan)
		}()
	}()

	return parsedDataChan
}

// pipelineGetHTML represents the step in the pipeline where raw HTML is gathered using POST requests, and emitted out using a channel.
func pipelineGetHTML[T any, V any](collector *colly.Collector, doneChan <-chan struct{}, data []V, recievedInfo IPipelineRecievedValue[V], formData PartialFormData, functions PipelineFunctions[T, V]) chan PipelineResponse[*goquery.Selection] {
	//Make channel for outputting raw HTML
	rawHTMLChan := make(chan PipelineResponse[*goquery.Selection])

	//Get HTML concurrently
	go func() {
		var wg sync.WaitGroup

		//If no data, return recieved info
		if len(data) == 0 {
			select {
			case rawHTMLChan <- PipelineResponse[*goquery.Selection]{Value: recievedInfo.Html(), Err: nil}:
			case <-doneChan:
			}
		}

		//Scrape in parallel
		for _, piece := range data {
			wg.Add(1)

			go func(piece V) {
				defer wg.Done()

				//Check if done's been called
				select {
				case <-doneChan:
					return
				default:
				}

				//Get HTML
				var html *goquery.Selection
				var err error

				if recievedInfo.Equal(piece) {
					html = recievedInfo.Html()
				} else {
					_, html, err = PostTo(collector, formData.Base, formData.Url, functions.GenFormData(functions.ToFormData(piece), formData))
				}

				//Try emitting HTML to channel
				select {
				case rawHTMLChan <- PipelineResponse[*goquery.Selection]{Value: html, Err: err}:
				case <-doneChan:
				}
			}(piece)
		}

		//Close channel once done
		go func() {
			wg.Wait()
			close(rawHTMLChan)
		}()
	}()

	return rawHTMLChan
}
