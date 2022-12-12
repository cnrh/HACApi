package utils

import (
	"errors"
	"strings"

	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

var ErrorInvalidCredentials = errors.New("invalid credentials")

// login logs a colly collector into Home Access Center.
func login(url, username, password string) (*colly.Collector, error) {
	// Get the base of the URL.
	base := strings.Split(url, "//")[1]

	// Create a new Colly collector.
	collector := colly.NewCollector(
		colly.AllowedDomains(base),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)

	// Create a channel to pass the request verification token into from HTML.
	reqVerChan := make(chan string, 1)

	// Create a channel to signal any errors.
	errChan := make(chan error, 1)

	// Retrieve request verification token from HTML.
	collector.OnHTML("input[name='__RequestVerificationToken']", func(elem *colly.HTMLElement) {
		reqVerChan <- elem.Attr("value")
	})

	// Handle any errors.
	collector.OnError(func(r *colly.Response, err error) {
		errChan <- err
	})

	// Form login URL.
	loginURL := url + repository.LOGIN_ROUTE

	// Get request verification token, abort if there are any errors.
	err := collector.Visit(loginURL)
	collector.Wait()

	if err != nil {
		return nil, err
	}

	collector = collector.Clone()

	// Get Request Verification Token or return any errors.
	var reqVerToken string

	select {
	// Request Verification Token obtained.
	case reqVerToken = <-reqVerChan:
	// Error obtained.
	case err := <-errChan:
		return nil, err
	}

	// Create payload data.
	payload := map[string]string{
		"__RequestVerificationToken": reqVerToken,
		"LogOnDetails.UserName":      username,
		"LogOnDetails.Password":      password,
		"SCKTY00328510CustomEnabled": "true",
		"SCKTY00436568CustomEnabled": "true",
		"Database":                   "10",
		"VerificationOption":         "UsernamePassword",
		"tempUN":                     "",
		"tempPW":                     "",
	}

	// Channel to signal whether the login failed or not.
	loginWrongChan := make(chan bool, 1)

	// Form the expected URL to be redirected to after the request.
	expectedURL := url + "/HomeAccess/Classes/Classwork"

	// Check if we are at expected URL. If not, the login has failed.
	collector.OnResponse(func(res *colly.Response) {
		if res.Request.URL.String() != expectedURL {
			loginWrongChan <- true
		}
	})

	// Set request headers.
	collector.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Host", base)
		req.Headers.Set("Origin", url)
		req.Headers.Set("Referer", base)
		req.Headers.Set("__RequestVerificationToken", reqVerToken)
	})

	// Handle errors.
	collector.OnError(func(r *colly.Response, err error) {
		errChan <- err
	})

	// Post to login.
	err = collector.Post(loginURL, payload)
	collector.Wait()

	// Check if login went through.
	if err != nil {
		return nil, err
	}

	// Handle any errors.
	select {
	// Credentials were wrong.
	case <-loginWrongChan:
		return nil, ErrorInvalidCredentials
	// Other error.
	case err := <-errChan:
		return nil, err
	default:
	}

	// Return logged-in collector.
	return collector, nil
}
