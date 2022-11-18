package utils

import (
	"strings"

	"github.com/Threqt1/HACApi/pkg/repository"
	"github.com/gocolly/colly"
)

// makePayloadString converts a map containing payload entires into a string payload for a Colly request.
func makePayloadString(payload *map[string]string) string {
	//Make the builder and allocate
	builder := strings.Builder{}
	builder.Grow(30000)

	//Convert the map into a payload string
	for key, val := range *payload {
		builder.WriteString(key)
		builder.WriteString("=")
		builder.WriteString(val)
		builder.WriteString("&")
	}

	//Return the payload string except the last extraenous &
	return strings.TrimSuffix(builder.String(), "&")
}

// Login logs a colly collector into home access center.
func Login(base, username, password string) *colly.Collector {
	//Create a new async collector
	collector := colly.NewCollector(
		colly.AllowedDomains(base),
		colly.Async(true),
	)

	//Create a channel to get the request verification token from HTML
	reqVerChan := make(chan string, 1)

	//Retrieve request verification token from HTML
	collector.OnHTML("input[name='__RequestVerificationToken']", func(elem *colly.HTMLElement) {
		reqVerChan <- elem.Attr("value")
	})

	//Form login URL
	loginURL := "https://" + base + repository.LOGIN_ROUTE

	//Get request verification token, abort if failed
	err := collector.Visit(loginURL)
	collector.Wait()

	if err != nil {
		return nil
	}

	//Create clone collector, let's login for real
	collector = collector.Clone()

	//Get Request Verification token
	reqVerToken := <-reqVerChan

	//Create payload data
	payloadData := map[string]string{
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

	//Make payload string
	payload := makePayloadString(&payloadData)

	//Channel to recieve whether login failed or not
	loginWrongChan := make(chan bool, 1)

	//Form URL we expect to be at after response
	expectedURL := "https://" + base + "/HomeAccess/Classes/Classwork"

	//Check if we are at expected URL. If not, login failed
	collector.OnResponse(func(res *colly.Response) {
		if res.Request.URL.String() != expectedURL {
			loginWrongChan <- true
		} else {
			loginWrongChan <- false
		}
	})

	//Set request headers
	collector.OnRequest(func(req *colly.Request) {
		req.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36")
		req.Headers.Set("Host", base)
		req.Headers.Set("Origin", "https://"+base)
		req.Headers.Set("Referer", base)
		req.Headers.Set("__RequestVerificationToken", reqVerToken)
	})

	//Post to login
	err = collector.PostRaw(loginURL, []byte(payload))
	collector.Wait()

	//Check if login went through
	if <-loginWrongChan || err != nil {
		return nil
	}

	//Return logged in collector
	return collector
}
