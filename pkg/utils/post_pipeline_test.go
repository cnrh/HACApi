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

// Represents a test input into a pipeline.
type testPipeline_Data struct {
	I int // The value the dummy input holds.
}

// Get HTML for a dummy input.
func (data testPipeline_Data) Html() *goquery.Selection {
	// Create a dummy GoQuery HTML page.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body><h1>` + strconv.Itoa(data.I) + `</h1><div class="fd">A,B,C,D,E</div></body></html>`))
	if err != nil {
		return nil
	}
	return doc.Find("body")
}

// Check equality between dummy inputs.
func (data testPipeline_Data) Equal(v int) bool {
	return data.I == v
}

// Represents a test output from pipeline.
type testPipeline_Return struct {
	J  int    // The value of the dummy input.
	FD string // The value of the form data carried with the input.
}

// Represents a dummy scraper.
type testPipeline_DummyScraper struct{}

// Represents the Login method for a dummy scraper (not needed).
func (scraper testPipeline_DummyScraper) Login(base, username, password string) (*colly.Collector, error) {
	return nil, nil
}

// Represents the Navigate method for a dummy scraper (not needed).
func (scraper testPipeline_DummyScraper) Navigate(collector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, nil
}

// Represents the Post method for a dummy scraper.
func (scraper testPipeline_DummyScraper) Post(collector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	// Embed the form data and the given I value into HTML, and return it.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(`<html><body><h1>` + formData["I"] + `</h1><div class="fd">` + formData["__VIEWSTATE"] + `,` + formData["__VIEWSTATEGENERATOR"] + `,` + formData["__EVENTVALIDATION"] + `,` + formData["__URL"] + `,` + formData["__BASE"] + `</div></body></html>`))
	if err != nil {
		return nil, nil, err
	}
	return nil, doc.Find("body"), nil
}

// Functions for a pipeline.
var testPipeline_Funcs = PipelineFunctions[testPipeline_Return, int]{
	GenFormData: func(s string, pfd PartialFormData) map[string]string {
		// Embed I into form data.
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
		// Convert the h1 value from the HTML to a number.
		num, _ := strconv.Atoi(s.Find("h1").First().Text())
		// Parse form data as well, and return both.
		fd := s.Find(".fd").Text()
		return testPipeline_Return{J: num, FD: fd}
	},
	// Convert int to string for form data embedding.
	ToFormData: strconv.Itoa,
}

// Static form data.
var testPipeline_Formdata = &PartialFormData{
	ViewState:       "A",
	ViewStateGen:    "B",
	EventValidation: "C",
	Url:             "D",
	Base:            "E",
}

// Test if GeneratePipeline() works with a recieved value only.
func TestGeneratePipeline_RecievedValueOnly(t *testing.T) {
	// Set up test pipeline data.
	recieved := testPipeline_Data{I: 1}
	data := []int{1}

	// Make expected value.
	expected := []testPipeline_Return{{J: 1, FD: "A,B,C,D,E"}}

	// Test.
	parsed, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)
	if err != nil {
		t.Fatalf("Failed for GeneratePipeline() Recieved Value Only:\n%v", err)
	}

	if diff := cmp.Diff(expected, parsed); diff != "" {
		t.Fatalf("Failed for GeneratePipeline() Recieved Value Only (-want, +got):\n%s", diff)
	}
}

// Test if GeneratePipeline() works with no data, but arecieved value.
func TestGeneratePipeline_NoValues_WithRecieved(t *testing.T) {
	// Set up test pipeline data.
	recieved := testPipeline_Data{I: 1}
	data := []int{}

	// Make expected value.
	expected := []testPipeline_Return{{J: 1, FD: "A,B,C,D,E"}}

	// Test.
	parsed, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)
	if err != nil {
		t.Fatalf("Failed for GeneratePipeline() No Values With Recieved:\n%v", err)
	}

	if diff := cmp.Diff(expected, parsed); diff != "" {
		t.Fatalf("Failed for GeneratePipeline() No Values With Recieved (-want, +got):\n%s", diff)
	}
}

// Test if GeneratePipeline() works with multiple values as well as a recieved value.
func TestGeneratePipeline_MultipleValues_WithRecieved(t *testing.T) {
	// Set up test pipeline data.
	recieved := testPipeline_Data{I: 1}
	data := []int{1, 2, 3, 4, 5}

	// Create a sorter to compare two slices.
	sorter := cmpopts.SortSlices(func(A, B testPipeline_Return) bool {
		return A.J < B.J
	})

	// Make expected value.
	expected := []testPipeline_Return{{J: 1, FD: "A,B,C,D,E"}, {J: 2, FD: "A,B,C,D,E"}, {J: 3, FD: "A,B,C,D,E"}, {J: 4, FD: "A,B,C,D,E"}, {J: 5, FD: "A,B,C,D,E"}}

	// Test.
	parsed, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)
	if err != nil {
		t.Fatalf("Failed for GeneratePipeline() Multiple Values With Recieved:\n%v", err)
	}

	if diff := cmp.Diff(expected, parsed, sorter); diff != "" {
		t.Fatalf("Failed for GeneratePipeline() Multiple Values With Recieved (-want, +got):\n%s", diff)
	}
}

// Test if GeneratePipeline() works with multiple values, but no recieved.
func TestGeneratePipeline_MultipleValues_WithoutRecieved(t *testing.T) {
	// Set up test pipeline data.
	recieved := testPipeline_Data{I: -1}
	data := []int{1, 2, 3, 4, 5}

	// Create a sorter to compare two slices.
	sorter := cmpopts.SortSlices(func(A, B testPipeline_Return) bool {
		return A.J < B.J
	})

	// Make expected value.
	expected := []testPipeline_Return{{J: 1, FD: "A,B,C,D,E"}, {J: 2, FD: "A,B,C,D,E"}, {J: 3, FD: "A,B,C,D,E"}, {J: 4, FD: "A,B,C,D,E"}, {J: 5, FD: "A,B,C,D,E"}}

	// Test.
	parsed, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)
	if err != nil {
		t.Fatalf("Failed for GeneratePipeline() Multiple Values Without Recieved:\n%v", err)
	}

	if diff := cmp.Diff(expected, parsed, sorter); diff != "" {
		t.Fatalf("Failed for GeneratePipeline() Multiple Values Without Recieved (-want, +got):\n%s", diff)
	}
}

// Represents a dummy scraper which always fails.
type testPipeline_DummyBadHTMLScraper struct{}

// Represents the Login method for a dummy scraper (not needed).
func (scraper testPipeline_DummyBadHTMLScraper) Login(base, username, password string) (*colly.Collector, error) {
	return nil, nil
}

// Represents the Navigate method for a dummy scraper (not needed).
func (scraper testPipeline_DummyBadHTMLScraper) Navigate(collector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, nil
}

// Represents the Post method for a dummy scraper, should always error out.
func (scraper testPipeline_DummyBadHTMLScraper) Post(collector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, ErrorBadHTML
}

// Test if GeneratePipeline() errors out if the POST request fails.
func TestGeneratePipeline_MalformedHTML(t *testing.T) {
	// Set up test pipeline data.
	recieved := testPipeline_Data{I: -1}
	data := []int{1, 2, 3, 4, 5}

	// Test.
	_, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyBadHTMLScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)

	if err == nil {
		t.Fatalf("Failed for GeneratePipeline() Malformed HTML (-want, +got):\n- %v\n+ nil", ErrorBadHTML)
	}

	if !errors.Is(err, ErrorBadHTML) {
		t.Fatalf("Failed for GeneratePipeline() Malformed HTML (-want, +got):\n- %v\n+ %v", ErrorBadHTML, err)
	}
}

// Represents a dummy scraper which always returns nil HTML.
type testPipeline_DummyNilHTMLScraper struct{}

// Represents the Login method for a dummy scraper (not needed).
func (scraper testPipeline_DummyNilHTMLScraper) Login(base, username, password string) (*colly.Collector, error) {
	return nil, nil
}

// Represents the Navigate method for a dummy scraper (not needed).
func (scraper testPipeline_DummyNilHTMLScraper) Navigate(collector *colly.Collector, base, url string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, nil
}

// Represents the Post method for a dummy scraper, should always return nil HTML.
func (scraper testPipeline_DummyNilHTMLScraper) Post(collector *colly.Collector, base, url string, formData map[string]string) (*colly.Collector, *goquery.Selection, error) {
	return nil, nil, nil
}

// Test if GeneratePipeline() errors out if the returned HTML is nil.
func TestGeneratePipeline_NilHTML(t *testing.T) {
	// Set up test pipeline data.
	recieved := testPipeline_Data{I: -1}
	data := []int{1, 2, 3, 4, 5}

	// Test.
	_, err := GeneratePipeline[testPipeline_Return, int](testPipeline_DummyNilHTMLScraper{}, nil, data, recieved, testPipeline_Formdata, testPipeline_Funcs)

	if err == nil {
		t.Fatalf("Failed for GeneratePipeline() Malformed HTML (-want, +got):\n- %v\n+ nil", ErrorBadHTML)
	}

	if !errors.Is(err, ErrorBadHTML) {
		t.Fatalf("Failed for GeneratePipeline() Malformed HTML (-want, +got):\n- %v\n+ %v", ErrorBadHTML, err)
	}
}
