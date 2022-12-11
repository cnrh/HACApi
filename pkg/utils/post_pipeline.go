package utils

import (
	"errors"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// pipelineResponse represents a intermediary
// response between pipelines
type pipelineResponse[T any] struct {
	Value T
	Err   error
}

// PartialFormData represents partial formdata
// for a POST request.
type PartialFormData struct {
	ViewState       string // viewstate formdata entry
	ViewStateGen    string // viewstategen formdata entry
	EventValidation string // eventvalidation formdata entry
	Url             string // url for the request
	Base            string // base url for the request
}

// PipelineFunctions describes the functions needed in order
// for the pipeline to correctly utilize the provided information.
type PipelineFunctions[T any, V any] struct {
	GenFormData func(string, PartialFormData) map[string]string
	Parse       func(*goquery.Selection) T
	ToFormData  func(V) string
}

// PipleineRecievedValue represents an interface which
// contains the needed methods for the pipeline to correctly
// handle the recievedInformation parameter.
type PipelineRecievedValue[V any] interface {
	Html() *goquery.Selection
	Equal(V) bool
}

// Represents an error due to no found HTML
var ErrorBadHTML = errors.New("bad html")

// GeneratePipeline creates a new pipeline which will gather data from a POST request, parse it, and return it in an array format. T represents the model struct,
// V represents the recieved value's type.
func GeneratePipeline[T any, V any](scraper repository.ScraperProvider, collector *colly.Collector, data []V, recievedInfo PipelineRecievedValue[V], formData PartialFormData, functions PipelineFunctions[T, V]) ([]T, error) {
	// Make a done channel for cancelling on error
	doneChan := make(chan struct{})
	defer close(doneChan)

	// Recieve parsed data
	parsedDataChan := pipelineParseHTML(scraper, collector, doneChan, data, recievedInfo, formData, functions)

	// Append recieved data to array
	dataArray := make([]T, 0, len(data))

	for res := range parsedDataChan {
		// If there's an error, cancel
		if res.Err != nil {
			return nil, res.Err
		}

		dataArray = append(dataArray, res.Value)
	}

	return dataArray, nil
}

// pipelineParseHTML represents the step in the pipeline where raw HTML is recieved through a channel, parsed, and emitted out through another channel.
func pipelineParseHTML[T any, V any](scraper repository.ScraperProvider, collector *colly.Collector, doneChan <-chan struct{}, data []V, recievedInfo PipelineRecievedValue[V], formData PartialFormData, functions PipelineFunctions[T, V]) chan pipelineResponse[T] {
	// Make a channel to emit parsed data/errors
	parsedDataChan := make(chan pipelineResponse[T])

	go func() {
		// Recieve raw HTML
		rawHTMLChan := pipelineGetHTML(scraper, collector, doneChan, data, recievedInfo, formData, functions)

		var wg sync.WaitGroup

		// Parse HTML concurrently
		for res := range rawHTMLChan {
			// If error, cascade it down and break
			if res.Err != nil {
				parsedDataChan <- pipelineResponse[T]{Err: res.Err}
				break
			}

			// If no HTML, error out
			if res.Value == nil {
				parsedDataChan <- pipelineResponse[T]{Err: ErrorBadHTML}
				break
			}

			// Otherwise, start goroutine to parse
			wg.Add(1)
			go func(res pipelineResponse[*goquery.Selection]) {
				defer wg.Done()

				// Check if done was called
				select {
				case <-doneChan:
					return
				default:
				}

				// Parse HTML
				parsedData := functions.Parse(res.Value)

				// Try emitting parsed data
				select {
				case parsedDataChan <- pipelineResponse[T]{Value: parsedData, Err: nil}:
				case <-doneChan:
				}
			}(res)
		}

		// Close channel once finished
		go func() {
			wg.Wait()
			close(parsedDataChan)
		}()
	}()

	return parsedDataChan
}

// pipelineGetHTML represents the step in the pipeline where raw HTML is gathered using POST requests, and emitted out using a channel.
func pipelineGetHTML[T any, V any](scraper repository.ScraperProvider, collector *colly.Collector, doneChan <-chan struct{}, data []V, recievedInfo PipelineRecievedValue[V], formData PartialFormData, functions PipelineFunctions[T, V]) chan pipelineResponse[*goquery.Selection] {
	// Make channel for outputting raw HTML
	rawHTMLChan := make(chan pipelineResponse[*goquery.Selection])

	// Get HTML concurrently
	go func() {
		var wg sync.WaitGroup

		// If no data, return recieved info
		if len(data) == 0 {
			select {
			case rawHTMLChan <- pipelineResponse[*goquery.Selection]{Value: recievedInfo.Html(), Err: nil}:
			case <-doneChan:
			}
		}

		// Scrape in parallel
		for _, piece := range data {
			wg.Add(1)

			go func(piece V) {
				defer wg.Done()

				// Check if done's been called
				select {
				case <-doneChan:
					return
				default:
				}

				// Get HTML
				var html *goquery.Selection
				var err error

				if recievedInfo.Equal(piece) {
					html = recievedInfo.Html()
				} else {
					_, html, err = scraper.Post(collector, formData.Base, formData.Url, functions.GenFormData(functions.ToFormData(piece), formData))
				}

				// Try emitting HTML to channel
				select {
				case rawHTMLChan <- pipelineResponse[*goquery.Selection]{Value: html, Err: err}:
				case <-doneChan:
				}
			}(piece)
		}

		// Close channel once done
		go func() {
			wg.Wait()
			close(rawHTMLChan)
		}()
	}()

	return rawHTMLChan
}
