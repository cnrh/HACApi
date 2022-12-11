package utils

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Represents a test input into pipeline
type testPipeline_Data struct {
	I int
}

// Get HTML for dummy input
func (data testPipeline_Data) Html() *goquery.Selection {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body><h1>` + strconv.Itoa(data.I) + `</h1><div class="fd">A,B,C,D,E</div></body></html>`))
	if err != nil {
		return nil
	}
	return doc.Find("body")
}

// Check equality between dummy inputs
func (data testPipeline_Data) Equal(v int) bool {
	return data.I == v
}

// Represents a test output from pipeline
type testPipeline_Return struct {
	J  int
	FD string
}

// Represents the dummy scraper
type testPipeline_DummyScraper struct{}

func (scraper testPipeline_DummyScraper) Login(base, username, password string) (*colly.Collector, error) {
	return nil, nil
}

func (scraper testPipeline_DummyScraper) Navigate(collector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, nil
}

func (scraper testPipeline_DummyScraper) Post(collector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body><h1>` + formData["I"] + `</h1><div class="fd">` + formData["__VIEWSTATE"] + `,` + formData["__VIEWSTATEGENERATOR"] + `,` + formData["__EVENTVALIDATION"] + `,` + formData["__URL"] + `,` + formData["__BASE"] + `</div></body></html>`))
	if err != nil {
		return nil, nil, err
	}
	return nil, doc.Find("body"), nil
}

// Functions for pipeline
var testPipeline_Funcs = PipelineFunctions[testPipeline_Return, int]{
	GenFormData: func(s string, pfd PartialFormData) map[string]string {
		return map[string]string{
			"__VIEWSTATE":          pfd.ViewState,
			"__VIEWSTATEGENERATOR": pfd.ViewStateGen,
			"__EVENTVALIDATION":    pfd.EventValidation,
			"__URL":                pfd.Url,
			"__BASE":               pfd.Base,
			"I":                    s,
		}
	},
	Parse: func(s *goquery.Selection) testPipeline_Return {
		num, _ := strconv.Atoi(s.Find("h1").First().Text())
		fd := s.Find(".fd").Text()
		return testPipeline_Return{J: num, FD: fd}
	},
	ToFormData: strconv.Itoa,
}

// Static formdata
var testPipeline_Formdata = PartialFormData{
	ViewState:       "A",
	ViewStateGen:    "B",
	EventValidation: "C",
	Url:             "D",
	Base:            "E",
}

// Check if pipeline works with recieved value only
func TestGeneratePipeline_RecievedValueOnly(t *testing.T) {
	// Set up test pipeline data
	recieved := testPipeline_Data{I: 1}
	data := []int{1}

	// Make expected
	expected := []testPipeline_Return{{J: 1, FD: "A,B,C,D,E"}}

	// Make pipeline
	parsed, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)
	if err != nil {
		t.Fatalf("Failed for GeneratePipeline() Recieved Value Only:\n%v", err)
	}

	if diff := cmp.Diff(expected, parsed); diff != "" {
		t.Fatalf("Failed for GeneratePipeline() Recieved Value Only (-want, +got):\n%s", diff)
	}
}

// Check if pipeline works with no values in data, but a recieved one provided
func TestGeneratePipeline_NoValues_WithRecieved(t *testing.T) {
	recieved := testPipeline_Data{I: 1}
	data := []int{}

	// Make expected
	expected := []testPipeline_Return{{J: 1, FD: "A,B,C,D,E"}}

	// Make pipeline
	parsed, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)
	if err != nil {
		t.Fatalf("Failed for GeneratePipeline() No Values With Recieved:\n%v", err)
	}

	if diff := cmp.Diff(expected, parsed); diff != "" {
		t.Fatalf("Failed for GeneratePipeline() No Values With Recieved (-want, +got):\n%s", diff)
	}
}

// Check if pipeline works with multiple values in data and a recieved one provided
func TestGeneratePipeline_MultipleValues_WithRecieved(t *testing.T) {
	// Set up test pipeline data
	recieved := testPipeline_Data{I: 1}
	data := []int{1, 2, 3, 4, 5}

	sorter := cmpopts.SortSlices(func(A, B testPipeline_Return) bool {
		return A.J < B.J
	})

	// Make expected
	expected := []testPipeline_Return{{J: 1, FD: "A,B,C,D,E"}, {J: 2, FD: "A,B,C,D,E"}, {J: 3, FD: "A,B,C,D,E"}, {J: 4, FD: "A,B,C,D,E"}, {J: 5, FD: "A,B,C,D,E"}}

	// Make pipeline
	parsed, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)
	if err != nil {
		t.Fatalf("Failed for GeneratePipeline() Multiple Values With Recieved:\n%v", err)
	}

	if diff := cmp.Diff(expected, parsed, sorter); diff != "" {
		t.Fatalf("Failed for GeneratePipeline() Multiple Values With Recieved (-want, +got):\n%s", diff)
	}
}

// Check if pipeline works with multiple values without recieved
func TestGeneratePipeline_MultipleValues_WithoutRecieved(t *testing.T) {
	// Set up test pipeline data
	recieved := testPipeline_Data{I: -1}
	data := []int{1, 2, 3, 4, 5}

	sorter := cmpopts.SortSlices(func(A, B testPipeline_Return) bool {
		return A.J < B.J
	})

	// Make expected
	expected := []testPipeline_Return{{J: 1, FD: "A,B,C,D,E"}, {J: 2, FD: "A,B,C,D,E"}, {J: 3, FD: "A,B,C,D,E"}, {J: 4, FD: "A,B,C,D,E"}, {J: 5, FD: "A,B,C,D,E"}}

	// Make pipeline
	parsed, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)
	if err != nil {
		t.Fatalf("Failed for GeneratePipeline() Multiple Values Without Recieved:\n%v", err)
	}

	if diff := cmp.Diff(expected, parsed, sorter); diff != "" {
		t.Fatalf("Failed for GeneratePipeline() Multiple Values Without Recieved (-want, +got):\n%s", diff)
	}
}

// Represents the dummy scraper which always fails
type testPipeline_DummyBadHTMLScraper struct{}

func (scraper testPipeline_DummyBadHTMLScraper) Login(base, username, password string) (*colly.Collector, error) {
	return nil, nil
}

func (scraper testPipeline_DummyBadHTMLScraper) Navigate(collector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, nil
}

func (scraper testPipeline_DummyBadHTMLScraper) Post(collector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, ErrorBadHTML
}

// Check if pipeline works with malformed HTML output
func TestGeneratePipeline_MalformedHTML(t *testing.T) {
	// Set up test pipeline data
	recieved := testPipeline_Data{I: -1}
	data := []int{1, 2, 3, 4, 5}

	_, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyBadHTMLScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)

	if err == nil {
		t.Fatalf("Failed for GeneratePipeline() Malformed HTML (-want, +got):\n- %v\n+ nil", ErrorBadHTML)
	}

	if !errors.Is(err, ErrorBadHTML) {
		t.Fatalf("Failed for GeneratePipeline() Malformed HTML (-want, +got):\n- %v\n+ %v", ErrorBadHTML, err)
	}
}

// Represents the dummy scraper which always returns nil html
type testPipeline_DummyNilHTMLScraper struct{}

func (scraper testPipeline_DummyNilHTMLScraper) Login(base, username, password string) (*colly.Collector, error) {
	return nil, nil
}

func (scraper testPipeline_DummyNilHTMLScraper) Navigate(collector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, nil
}

func (scraper testPipeline_DummyNilHTMLScraper) Post(collector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, nil
}

// Check if pipeline works with nil HTML values
func TestGeneratePipeline_NilHTML(t *testing.T) {
	// Set up test pipeline data
	recieved := testPipeline_Data{I: -1}
	data := []int{1, 2, 3, 4, 5}

	_, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyNilHTMLScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)

	if err == nil {
		t.Fatalf("Failed for GeneratePipeline() Malformed HTML (-want, +got):\n- %v\n+ nil", ErrorBadHTML)
	}

	if !errors.Is(err, ErrorBadHTML) {
		t.Fatalf("Failed for GeneratePipeline() Malformed HTML (-want, +got):\n- %v\n+ %v", ErrorBadHTML, err)
	}
}
