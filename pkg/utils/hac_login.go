package utils

import (
	"errors"

	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// Login logs a colly collector into home access center.
func Login(base, username, password string) (*colly.Collector, error) {
	// Create a new async collector
	collector := colly.NewCollector(
		colly.AllowedDomains(base),
		colly.Async(true),
		colly.AllowURLRevisit(),
	)

	// Create a channel to get the request verification token from HTML
	reqVerChan := make(chan string, 1)

	// Retrieve request verification token from HTML
	collector.OnHTML("input[name='__RequestVerificationToken']", func(elem *colly.HTMLElement) {
		reqVerChan <- elem.Attr("value")
	})

	// Form login URL
	loginURL := "https://" + base + repository.LOGIN_ROUTE

	// Get request verification token, abort if failed
	err := collector.Visit(loginURL)
	collector.Wait()

	if err != nil {
		return nil, err
	}

	// Create clone collector, let's login for real
	collector = collector.Clone()

	// Get Request Verification token
	reqVerToken := <-reqVerChan

	// Create payload data
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

	// Channel to recieve whether login failed or not
	loginWrongChan := make(chan bool, 1)

	// Form URL we expect to be at after response
	expectedURL := "https://" + base + "/HomeAccess/Classes/Classwork"

	// Check if we are at expected URL. If not, login failed
	collector.OnResponse(func(res *colly.Response) {
		if res.Request.URL.String() != expectedURL {
			loginWrongChan <- true
		} else {
			loginWrongChan <- false
		}
	})

	// Set request headers
	collector.OnRequest(func(req *colly.Request) {
		req.Headers.Set("Host", base)
		req.Headers.Set("Origin", "https://"+base)
		req.Headers.Set("Referer", base)
		req.Headers.Set("__RequestVerificationToken", reqVerToken)
	})

	// Post to login
	err = collector.Post(loginURL, payload)
	collector.Wait()

	// Check if login went through
	if err != nil {
		return nil, err
	}

	if <-loginWrongChan {
		return nil, errors.New("invalid credentials")
	}

	// Return logged in collector
	return collector, nil
}
