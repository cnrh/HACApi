package parsers

import (
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
)

// parseClasswork takes in the initial page HTML, and outputs
// the parsed classwork.
func parseClasswork(html *goquery.Selection) models.Classwork {
	// Make a struct to store parsed classwork in, allocate if necessary
	classwork := models.Classwork{}

	// Get all the classes on the page
	classEles := html.Find(".AssignmentClass")

	// Allocate memory for the slice
	classwork.Entries = make([]models.ClassworkEntry, 0, classEles.Length())

	// Find marking period, try to make it into an int
	MarkingPerStr := strings.TrimSpace(html.Find("#plnMain_ddlReportCardRuns > option[selected='selected']").Text())
	MarkingPer, err := strconv.Atoi(MarkingPerStr)
	if err == nil {
		classwork.MarkingPeriod = MarkingPer
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Go through each class, parsing all assignments for each class and other data into a ClassworkEntry struct and
	// pushing it to the slice
	classEles.Each(func(classPos int, classEle *goquery.Selection) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			// Get classwork entry, push it to slice
			classworkEntry := parseClassworkEntry(classEle, classPos)
			mutex.Lock()
			defer mutex.Unlock()
			classwork.Entries = append(classwork.Entries, classworkEntry)
		}()
	})

	wg.Wait()

	return classwork
}

// parseClassworkEntry parses an individual class entry in the page. Meant for concurrency.
func parseClassworkEntry(classEle *goquery.Selection, classPos int) models.ClassworkEntry {
	// Create the entry
	classworkEntry := models.ClassworkEntry{}

	classworkEntry.Position = classPos

	// Split element title by space
	fullElementTitle := strings.TrimSpace(classEle.Find("a.sg-header-heading").First().Text())
	splitElementTitle := strings.Split(fullElementTitle, " ")

	// Get relevant information
	className := strings.TrimSpace(strings.Join(splitElementTitle[3:], " "))
	courseName := strings.TrimSpace(strings.Join(splitElementTitle[0:3], " "))

	// Assign it to object
	classworkEntry.Class.Name = className
	classworkEntry.Class.Course = courseName

	// Get average grade
	splitAverageText := strings.Split(strings.TrimSpace(classEle.Find("span.sg-header-heading").First().Text()), " ")
	classworkEntry.Average = strings.TrimSpace(splitAverageText[len(splitAverageText)-1])

	// Get all assignments
	assignments := classEle.Find("table.sg-asp-table:first-child tr.sg-asp-table-data-row")

	// Allocate space for assignments array
	classworkEntry.Assignments = make([]models.Assignment, 0, assignments.Length())

	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Loop through each assignment, parsing them
	assignments.Each(func(_ int, assignmentEle *goquery.Selection) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			assignment := parseClassworkAssignment(assignmentEle)
			mutex.Lock()
			defer mutex.Unlock()
			classworkEntry.Assignments = append(classworkEntry.Assignments, assignment)
		}()
	})

	wg.Wait()

	return classworkEntry
}

// parseClassworkAssignment parses an individual assignment inside a class cluster.
// Meant for concurrency.
func parseClassworkAssignment(assignmentEle *goquery.Selection) models.Assignment {
	// Create the assignment
	assignment := models.Assignment{}

	// Go through each td in the assignment's HTML row, using index to figure out
	// what data it represents
	assignmentEle.Find("td").Each(func(i int, dataEle *goquery.Selection) {
		// Parse text, return if there is none
		text := strings.TrimSpace(dataEle.Text())
		if text == "" {
			return
		}

		// Fill in data using i
		switch i {
		case 0:
			assignment.DueDate = text
		case 1:
			assignment.AssignedDate = text
		case 2:
			text = strings.TrimSpace(strings.Replace(text, "*", "", 1))
			assignment.Name = text
		case 3:
			assignment.Category = text
		case 4:
			assignment.Grade = text
			style, exists := dataEle.Attr("style")
			if exists && strings.Contains(style, "line-through") || text == "X" {
				assignment.Dropped = true
			}
		case 5:
			assignment.TotalPoints = text
		}
	})

	return assignment
}
