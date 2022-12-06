package parsers

import (
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
)

// parseIPR takes in the initial page HTMl, and
// returns the parsed IPR.
func parseIPR(html *goquery.Selection) models.IPR {
	// Make a struct to store parsed IPR info to
	ipr := models.IPR{}

	// Get date
	dateText := html.Find("#plnMain_ddlIPRDates > option[selected='selected']").Text()
	ipr.Date = strings.TrimSpace(dateText)

	// Get all IPR class HTML rows
	classEles := html.Find("table.sg-asp-table:first-child tr.sg-asp-table-data-row")

	// Allocate memory for the array
	ipr.Entries = make([]models.IPREntry, 0, classEles.Length())

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Parse each row
	classEles.Each(func(_ int, iprRowEle *goquery.Selection) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			iprEntry := parseIPREntry(iprRowEle)
			mutex.Lock()
			defer mutex.Unlock()
			ipr.Entries = append(ipr.Entries, iprEntry)
		}()
	})

	wg.Wait()

	return ipr
}

// parseIPREntry takes in a HTML selection representing the
// row which contains information about a specific class's
// progress report, and returns the parsed IPREntry for it.
func parseIPREntry(iprRowEle *goquery.Selection) models.IPREntry {
	// Create the entry
	iprEntry := models.IPREntry{}

	// Go through tds to parse relevant info
	iprRowEle.Find("td").Each(func(i int, dataEle *goquery.Selection) {
		// Parse text, return if there is none
		text := strings.TrimSpace(dataEle.Text())
		if len(text) == 0 {
			return
		}

		// Fill in data using i
		switch i {
		case 0:
			iprEntry.Class.Course = text
		case 1:
			iprEntry.Class.Name = text
		case 2:
			iprEntry.Class.Period = text
		case 3:
			iprEntry.Class.Teacher = text
		case 4:
			iprEntry.Class.Room = text
		case 5:
			iprEntry.Grade = text
		}
	})

	return iprEntry
}
