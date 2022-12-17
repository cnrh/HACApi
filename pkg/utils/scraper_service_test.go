package utils

import (
	"strings"
	"testing"

	"github.com/gocolly/colly"
	"github.com/google/go-cmp/cmp"
)

// scraperTest represents the expected value from NewScraper().
type scraperTest struct {
	Value *Scraper
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
	ts := CreateTestingServer()
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
	ts := CreateTestingServer()
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
	ts := CreateTestingServer()
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
	ts := CreateTestingServer()
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
	ts := CreateTestingServer()
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
	ts := CreateTestingServer()
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
	ts := CreateTestingServer()
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
	ts := CreateTestingServer()
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
