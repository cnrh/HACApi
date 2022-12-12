package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/gocolly/colly"
	"github.com/google/go-cmp/cmp"
)

// scraperTest represents the expected value from NewScraper().
type scraperTest struct {
	Value *Scraper
}

// parseRequestBody parses a http request body into a map.
func parseRequestBody(input string) map[string]string {
	splitParams := strings.Split(input, "&")
	output := make(map[string]string, len(splitParams))

	for _, param := range splitParams {
		splitParam := strings.Split(param, "=")
		output[splitParam[0]] = splitParam[1]
	}

	return output
}

// Create a test scraping server.
func createTestingServer() *httptest.Server {
	mux := http.NewServeMux()

	// Handle calls to repository.LOGIN_ROUTE.
	mux.HandleFunc("/HomeAccess/Account/LogOn", func(w http.ResponseWriter, r *http.Request) {
		// Get the static Login page.
		html, err := os.ReadFile("../../test/login.html")

		// Handle the POST request.
		if r.Method == "POST" {
			expected := map[string]string{
				"__RequestVerificationToken": "ABCD12345",
				"LogOnDetails.UserName":      "ABC",
				"LogOnDetails.Password":      "123",
				"SCKTY00328510CustomEnabled": "true",
				"SCKTY00436568CustomEnabled": "true",
				"Database":                   "10",
				"VerificationOption":         "UsernamePassword",
				"tempUN":                     "",
				"tempPW":                     "",
			}
			// Get the request body params.
			got, err := io.ReadAll(r.Body)

			// Compare them, if possible.
			if err == nil && reflect.DeepEqual(expected, parseRequestBody(string(got))) {
				// Redirect to classwork endpoint.
				http.Redirect(w, r, "/HomeAccess/Classes/Classwork", http.StatusSeeOther)
			} else {
				// Send back login HTML.
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(html))
			}
		}

		if err == nil {
			// Send login HTML if not a POST request.
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
		}
	})

	// Dummy response, just for handling redirect.
	mux.HandleFunc("/HomeAccess/Classes/Classwork", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<!doctype html><html></html>`))
	})

	// Handle a default static HTML page.
	mux.HandleFunc("/default", func(w http.ResponseWriter, r *http.Request) {
		// Send the default page.
		html, err := os.ReadFile("../../test/default.html")

		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
		}
	})

	// Dummy handler to test redirects.
	mux.HandleFunc("/redirect", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/HomeAccess/Classes/Classwork", http.StatusSeeOther)
	})

	// Handle a default post request.
	mux.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			expected := map[string]string{
				"A": "1",
				"B": "2",
				"C": "3",
			}
			// Read request body.
			got, err := io.ReadAll(r.Body)
			if err == nil && reflect.DeepEqual(expected, parseRequestBody(string(got))) {
				html, err := os.ReadFile("../../test/default.html")

				if err == nil {
					w.Header().Set("Content-Type", "text/html")
					w.Write([]byte(html))
				}
			} else {
				// Otherwise, error out.
				http.Redirect(w, r, "/error", http.StatusSeeOther)
			}
		}
	})

	// Redirect to handle a missing endpoint/malformed request.
	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		// Send static error page.
		html, err := os.ReadFile("../../test/error.html")

		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
		}
	})

	return httptest.NewServer(mux)
}

// Test creating a new scraper.
func TestNewScraper(t *testing.T) {
	// Expected value.
	expected := scraperTest{
		Value: &Scraper{},
	}

	// Test.
	scraper := NewScraper()

	if diff := cmp.Diff(expected, scraperTest{Value: scraper}); diff != "" {
		t.Fatalf("Failed for NewScraper() (-want, +got):\n%s", diff)
	}
}

// Test if Login() works with valid credentials.
func TestLogin_WithValidCredentials(t *testing.T) {
	// Create testing server and scraper.
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	// Test.
	collector, err := scraper.Login(ts.URL, "ABC", "123")

	if err != nil || collector == nil {
		t.Fatalf("Failed for Login() with valid credentials:\n%v", err)
	}
}

// Test if Login() errors out with invalid credentials.
func TestLogin_WithInvalidCredentials(t *testing.T) {
	// Create testing server and scraper.
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	// Test.
	collector, err := scraper.Login(ts.URL, "123", "ABC")

	if err != ErrorInvalidCredentials || collector != nil {
		t.Fatalf("Failed for Login() with invalid credentials")
	}
}

// Test if Login() errors out with an invalid URL.
func TestLogin_WithInvalidURL(t *testing.T) {
	// Create scraper.
	scraper := NewScraper()

	// Test.
	collector, err := scraper.Login("https://fake.url", "123", "ABC")

	if err == nil || collector != nil {
		t.Fatalf("Failed for Login() with invalid URL")
	}
}

// Test if Navigate() works with a valid URL.
func TestNavigate_WithValidURL(t *testing.T) {
	// Create testing server and scraper.
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	// Test.
	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())

	_, html, err := scraper.Navigate(initialCollector, ts.URL, "/default")

	if err != nil {
		t.Fatalf("Failed for Navigate() with valid url:\n%v", err)
	}

	// Confirm the correct page was returned.
	pageType, _ := html.Find("input[name='type']").First().Attr("value")

	if diff := cmp.Diff("default", pageType); diff != "" {
		t.Fatalf("Failed for Navigate() with valid url (-want, +got):\n%s", diff)
	}
}

// Test if Navigate() errors out with an invalid URL.
func TestNavigate_WithInvalidURL(t *testing.T) {
	// Create testing server and scraper.
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	// Test.
	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())

	_, _, err := scraper.Navigate(initialCollector, ts.URL, "/invalid")

	if err == nil {
		t.Fatalf("Failed for Navigate() with invalid url")
	}
}

// Test if Navigate() errors out if redirected.
func TestNavigate_WithRedirectBack(t *testing.T) {
	// Create testing server and scraper.
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	// Test.
	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())

	_, _, err := scraper.Navigate(initialCollector, ts.URL, "/redirect")

	if err == nil {
		t.Fatalf("Failed for Navigate() with redirect back")
	}
}

// Test if Post() works with valid form data and URL.
func TestPost_WithValidFormdata(t *testing.T) {
	// Create testing server and scraper.
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	// Create dummy collector and form data.
	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())
	formData := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}

	// Test.
	_, html, err := scraper.Post(initialCollector, ts.URL, "/post", formData)

	if err != nil {
		t.Fatalf("Failed for TestPost() with valid form data:\n%v", err)
	}

	// Confirm correct page was returned.
	pageType, _ := html.Find("input[name='type']").First().Attr("value")

	if diff := cmp.Diff("default", pageType); diff != "" {
		t.Fatalf("Failed for Post() with valid form data (-want, +got):\n%s", diff)
	}
}

// Test if Post() errors out with invalid form data.
func TestPost_WithInvalidFormData(t *testing.T) {
	// Create testing server and scraper.
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	// Create dummy collector and form data.
	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())
	formData := map[string]string{
		"A": "3",
		"B": "2",
		"C": "1",
	}

	_, _, err := scraper.Post(initialCollector, ts.URL, "/post", formData)

	if err == nil {
		t.Fatalf("Failed for Post() with invalid form data:\n%v", err)
	}
}

// Test if Post() errors out with an invalid URL.
func TestPost_WithInvalidURL(t *testing.T) {
	// Create testing server and scraper.
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	// Create dummy collector and form data.
	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())
	formData := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}

	_, _, err := scraper.Post(initialCollector, ts.URL, "/invalid", formData)

	if err == nil {
		t.Fatalf("Failed for Post() with an invalid URL:\n%v", err)
	}
}
