package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gocolly/colly"
	"github.com/google/go-cmp/cmp"
)

type scraperTest struct {
	Value *Scraper
}

// Create a test scraping server
func createTestingServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/HomeAccess/Account/LogOn", func(w http.ResponseWriter, r *http.Request) {
		html, err := os.ReadFile("../../test/login.html")

		if r.Method == "POST" {
			expecteddata := "Database=10&LogOnDetails.Password=123&LogOnDetails.UserName=ABC&SCKTY00328510CustomEnabled=true&SCKTY00436568CustomEnabled=true&VerificationOption=UsernamePassword&__RequestVerificationToken=ABCD12345&tempPW=&tempUN="
			data, err := io.ReadAll(r.Body)
			if err == nil && string(data) == expecteddata {
				http.Redirect(w, r, "/HomeAccess/Classes/Classwork", http.StatusSeeOther)
			} else {
				w.Header().Set("Content-Type", "text/html")
				w.Write([]byte(html))
			}
		}

		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
		}
	})

	mux.HandleFunc("/HomeAccess/Classes/Classwork", func(w http.ResponseWriter, r *http.Request) {
		html, err := os.ReadFile("../../test/login.html")

		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
		}
	})

	mux.HandleFunc("/existingendpoint", func(w http.ResponseWriter, r *http.Request) {
		html, err := os.ReadFile("../../test/existingendpoint.html")

		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
		}
	})

	mux.HandleFunc("/redirectendpoint", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/HomeAccess/Classes/Classwork", http.StatusSeeOther)
	})

	mux.HandleFunc("/postendpoint", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			expecteddata := "A=1&B=2&C=3"
			data, err := io.ReadAll(r.Body)
			if err == nil && string(data) == expecteddata {
				html, err := os.ReadFile("../../test/post.html")

				if err == nil {
					w.Header().Set("Content-Type", "text/html")
					w.Write([]byte(html))
				}
			} else {
				http.Redirect(w, r, "/error", http.StatusSeeOther)
			}
		}
	})

	mux.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		html, err := os.ReadFile("../../test/error.html")

		if err == nil {
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte(html))
		}
	})

	return httptest.NewServer(mux)
}

// Test creating a new scraper
func TestNewScraper(t *testing.T) {
	// Expected
	expected := scraperTest{
		Value: &Scraper{},
	}

	// Test
	scraper := NewScraper()

	if diff := cmp.Diff(expected, scraperTest{Value: scraper}); diff != "" {
		t.Fatalf("Failed for NewScraper() (-want, +got):\n%s", diff)
	}
}

func TestLogin_WithValidCredentials(t *testing.T) {
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	collector, err := scraper.Login(ts.URL, "ABC", "123")

	if err != nil || collector == nil {
		t.Fatalf("Failed for TestLogin() with valid credentials:\n%v", err)
	}
}

func TestLogin_WithInvalidCredentials(t *testing.T) {
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	collector, err := scraper.Login(ts.URL, "123", "ABC")

	if err != ErrorInvalidCredentials || collector != nil {
		t.Fatalf("Failed for TestLogin() with invalid credentials")
	}
}

func TestLogin_WithInvalidURL(t *testing.T) {
	scraper := NewScraper()

	collector, err := scraper.Login("https://fake.url", "123", "ABC")

	if err == nil || collector != nil {
		t.Fatalf("Failed for TestLogin() with invalid URL")
	}
}

func TestNavigate_WithValidURL(t *testing.T) {
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())

	_, html, err := scraper.Navigate(initialCollector, ts.URL, "/existingendpoint")

	if err != nil {
		t.Fatalf("Failed for TestNavigate() with valid url:\n%v", err)
	}

	if diff := cmp.Diff("Existing Endpoint", html.Find("h1").First().Text()); diff != "" {
		t.Fatalf("Failed for TestNavigate() with valid url (-want, +got):\n%s", diff)
	}
}

func TestNavigate_WithInvalidURL(t *testing.T) {
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())

	_, _, err := scraper.Navigate(initialCollector, ts.URL, "/nonexistantendpoint")

	if err == nil {
		t.Fatalf("Failed for TestNavigate() with invalid url")
	}
}

func TestNavigate_WithRedirectBack(t *testing.T) {
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())

	_, _, err := scraper.Navigate(initialCollector, ts.URL, "/redirectendpoint")

	if err == nil {
		t.Fatalf("Failed for TestNavigate() with redirect back")
	}
}

func TestPost_WithValidFormdata(t *testing.T) {
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())
	formdata := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}

	_, html, err := scraper.Post(initialCollector, ts.URL, "/postendpoint", formdata)

	if err != nil {
		t.Fatalf("Failed for TestPost() with valid form data:\n%v", err)
	}

	if diff := cmp.Diff("Post Successful", html.Find("h1").First().Text()); diff != "" {
		t.Fatalf("Failed for TestPost() with valid form data (-want, +got):\n%s", diff)
	}
}

func TestPost_WithInvalidFormData(t *testing.T) {
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())
	formdata := map[string]string{
		"A": "3",
		"B": "2",
		"C": "1",
	}

	_, _, err := scraper.Post(initialCollector, ts.URL, "/postendpoint", formdata)

	if err == nil {
		t.Fatalf("Failed for TestPost() with invalid form data:\n%v", err)
	}
}

func TestPost_WithInvalidURL(t *testing.T) {
	ts := createTestingServer()
	defer ts.Close()

	scraper := NewScraper()

	initialCollector := colly.NewCollector(colly.Async(true), colly.AllowedDomains(strings.Split(ts.URL, "//")[1]), colly.AllowURLRevisit())
	formdata := map[string]string{
		"A": "1",
		"B": "2",
		"C": "3",
	}

	_, _, err := scraper.Post(initialCollector, ts.URL, "/nonexistantendpoint", formdata)

	if err == nil {
		t.Fatalf("Failed for TestPost() with invalid form data:\n%v", err)
	}
}
