package parsers

import (
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
)

// parseReportCard takes in raw HTML and parses it into a report
// card model.
func parseReportCard(html *goquery.Selection) models.ReportCard {
	// Make struct to store parsed report card
	reportCard := models.ReportCard{}

	// Get all entries
	reportCardEntryEles := html.Find("tr.sg-asp-table-data-row")

	// Allocate space for array
	reportCard.Entries = make([]models.ReportCardEntry, 0, reportCardEntryEles.Length())

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Go through each entry to parse it
	reportCardEntryEles.Each(func(_ int, reportCardEntryEle *goquery.Selection) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			reportCardEntry := parseReportCardEntry(reportCardEntryEle)
			mutex.Lock()
			defer mutex.Unlock()
			reportCard.Entries = append(reportCard.Entries, reportCardEntry)
		}()
	})

	wg.Wait()

	return reportCard
}

// parseReportCardEntry takes in a report card entry HTML element and parses it
// into a ReportCardEntry struct.
func parseReportCardEntry(reportCardEntryEle *goquery.Selection) models.ReportCardEntry {
	// Make a struct to store parsed data
	reportCardEntry := models.ReportCardEntry{}

	// Go through each td, using i to match text to the corresponding field
	reportCardEntryEle.Find("td").Each(func(i int, dataEle *goquery.Selection) {
		// Parse text, return if there is none
		text := strings.TrimSpace(dataEle.Text())
		if text == "" {
			return
		}

		// Fill in data using i
		switch i {
		case 0:
			reportCardEntry.Class.Course = text
		case 1:
			reportCardEntry.Class.Name = text
		case 2:
			reportCardEntry.Class.Period = text
		case 3:
			reportCardEntry.Class.Teacher = text
		case 4:
			reportCardEntry.Class.Room = text
		case 5:
			reportCardEntry.AttemptedCredit = text
		case 6:
			reportCardEntry.EarnedCredit = text
		case 7:
			reportCardEntry.Averages.First = text
		case 8:
			reportCardEntry.Averages.Second = text
		case 9:
			reportCardEntry.Averages.Third = text
		case 10:
			reportCardEntry.Averages.Exam1 = text
		case 11:
			reportCardEntry.Averages.Sem1 = text
		case 12:
			reportCardEntry.Averages.Fourth = text
		case 13:
			reportCardEntry.Averages.Fifth = text
		case 14:
			reportCardEntry.Averages.Sixth = text
		case 15:
			reportCardEntry.Averages.Exam2 = text
		case 16:
			reportCardEntry.Averages.Sem2 = text
		case 17:
			reportCardEntry.Conduct.First = text
		case 18:
			reportCardEntry.Conduct.Second = text
		case 19:
			reportCardEntry.Conduct.Third = text
		case 20:
			reportCardEntry.Conduct.Fourth = text
		case 21:
			reportCardEntry.Conduct.Fifth = text
		case 22:
			reportCardEntry.Conduct.Sixth = text
		case 23:
			reportCardEntry.Comments.First = text
		case 24:
			reportCardEntry.Comments.Second = text
		case 25:
			reportCardEntry.Comments.Third = text
		case 26:
			reportCardEntry.Comments.Fourth = text
		case 27:
			reportCardEntry.Comments.Fifth = text
		case 28:
			reportCardEntry.Comments.Sixth = text
		case 29:
			reportCardEntry.Absences.ExcusedAbsence = text
		case 30:
			reportCardEntry.Absences.UnexcusedAbsence = text
		case 31:
			reportCardEntry.Absences.ExcusedTardy = text
		case 32:
			reportCardEntry.Absences.UnexcusedTardy = text
		}
	})

	return reportCardEntry
}
