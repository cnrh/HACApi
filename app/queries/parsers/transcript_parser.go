package parsers

import (
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
)

// parseTranscript takes in raw HTML and parses it into
// a transcript struct.
func parseTranscript(html *goquery.Selection) models.Transcript {
	// Create struct to hold parsed transcript
	transcript := models.Transcript{}

	// Find all group HTML elements
	transcriptGroupEles := html.Find("td.sg-transcript-group")

	// Allocate memory for the slice
	transcript.Entries = make([]models.TranscriptGroup, 0, transcriptGroupEles.Length())

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Parse each semester entry
	transcriptGroupEles.Each(func(transcriptGroupPos int, transcriptGroupEle *goquery.Selection) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			transcriptGroup := parseTranscriptGroup(transcriptGroupEle)
			mutex.Lock()
			defer mutex.Unlock()
			// Append group to slice
			transcript.Entries = append(transcript.Entries, transcriptGroup)
		}()
	})

	// Parse GPA concurrently
	html.Find(".sg-content-grid > table > tbody > tr:nth-last-child(2) > td > table > tbody > tr.sg-asp-table-data-row").Each(func(gpaType int, transcriptGPAEle *goquery.Selection) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			transcriptGPA := parseTranscriptGPA(transcriptGPAEle)
			mutex.Lock()
			defer mutex.Unlock()
			// Assign GPA based on its index (type)
			switch gpaType {
			case 0:
				transcript.Weighted = transcriptGPA
			case 1:
				transcript.Unweighted = transcriptGPA
			}
		}()
	})

	wg.Wait()

	return transcript
}

// parseTranscriptGroup takes in a HTML element representing a transcript group,
// and parses it into a TranscriptGroup struct.
func parseTranscriptGroup(transcriptGroupEle *goquery.Selection) models.TranscriptGroup {
	// Make struct to put parsed data into
	transcriptGroup := models.TranscriptGroup{}

	// Get all entries in the group element
	transcriptGroupEntryEles := transcriptGroupEle.Find("table.sg-asp-table tr.sg-asp-table-data-row")

	// Allocate memory
	transcriptGroup.Entries = make([]models.TranscriptGroupEntry, 0, transcriptGroupEntryEles.Length())

	// Parse the top table for information about the group
	transcriptGroupEle.Find("table:first-child td:nth-child(even)").Each(func(i int, dataEle *goquery.Selection) {
		// Parse text, return if there is none
		text := strings.TrimSpace(dataEle.Text())
		if text == "" {
			return
		}

		switch i {
		case 0:
			transcriptGroup.Year = text
		case 1:
			transcriptGroup.Semester = text
		case 2:
			transcriptGroup.GradeLevel = text
		case 3:
			transcriptGroup.Building = text
		}
	})

	// Get total credit
	totalCreditText := transcriptGroupEle.Find("table[style='float:right'] td:last-child").Text()
	transcriptGroup.TotalCredit = strings.TrimSpace(totalCreditText)

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Parse each entry in the table
	transcriptGroupEntryEles.Each(func(_ int, transcriptGroupEntryEle *goquery.Selection) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			transcriptGroupEntry := parseTranscriptGroupEntry(transcriptGroupEntryEle)
			mutex.Lock()
			defer mutex.Unlock()
			transcriptGroup.Entries = append(transcriptGroup.Entries, transcriptGroupEntry)
		}()
	})

	wg.Wait()

	return transcriptGroup
}

// parseTranscriptGroupEntry takes in a HTML element representing a transcript group entry, and
// parses it to return a TranscriptGroupEntry struct.
func parseTranscriptGroupEntry(transcriptGroupEntryEle *goquery.Selection) models.TranscriptGroupEntry {
	// Make struct to put parsed data into
	transcriptGroupEntry := models.TranscriptGroupEntry{}

	// Parse each td, using i to match the text to what it describes
	transcriptGroupEntryEle.Find("td").Each(func(i int, dataEle *goquery.Selection) {
		// Parse text, return if there is none
		text := strings.TrimSpace(dataEle.Text())
		if text == "" {
			return
		}

		switch i {
		case 0:
			transcriptGroupEntry.Class.Course = text
		case 1:
			transcriptGroupEntry.Class.Name = text
		case 2:
			transcriptGroupEntry.Average = text
		case 3:
			transcriptGroupEntry.Credit = text
		}
	})

	return transcriptGroupEntry
}

// parseTranscriptGPA takes in a HTML object representing a GPA record, and
// parses it to return a TranscriptGPA struct.
func parseTranscriptGPA(transcriptGPAEle *goquery.Selection) models.TranscriptGPA {
	// Struct for parsed data
	transcriptGPA := models.TranscriptGPA{}

	// Loop through all TDs, fill in necessary data
	transcriptGPAEle.Find("td").Each(func(i int, dataEle *goquery.Selection) {
		// Parse text, return if there is none
		text := strings.TrimSpace(dataEle.Text())
		if text == "" {
			return
		}

		switch i {
		case 0:
			gpaType := strings.TrimSuffix(text, "*")
			transcriptGPA.Type = gpaType
		case 1:
			transcriptGPA.GPA = text
		case 2:
			transcriptGPA.Rank = text
		case 3:
			transcriptGPA.Quartile = text
		}
	})

	return transcriptGPA
}
