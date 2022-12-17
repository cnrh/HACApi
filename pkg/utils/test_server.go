package utils

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
)

// ExpectedServerResponse represents a response to a
// request, which can be used for comparing expected vs
// actual.
type ExpectedServerResponse[V interface{}] struct {
	Status int
	Body   V
}

// ParseRequestBody parses a http request body into a map.
func ParseRequestBody(input string) map[string]string {
	splitParams := strings.Split(input, "&")
	output := make(map[string]string, len(splitParams))

	for _, param := range splitParams {
		splitParam := strings.Split(param, "=")
		output[splitParam[0]] = splitParam[1]
	}

	return output
}

// Create a test scraping server.
func CreateTestingServer() *httptest.Server {
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
			if err == nil && reflect.DeepEqual(expected, ParseRequestBody(string(got))) {
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
			if err == nil && reflect.DeepEqual(expected, ParseRequestBody(string(got))) {
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
