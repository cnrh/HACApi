package parsers

import (
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/Threqt1/HACApi/app/models"
)

// ParseSchedule takes in raw HTML and parses it into a schedule
// model.
func ParseSchedule(html *goquery.Selection) models.Schedule {
	//Make a struct to store parsed schedule
	schedule := models.Schedule{}

	//Get all schedule entry HTML elements
	scheduleEntryEles := html.Find("tr.sg-asp-table-data-row")

	//Allocate memory for array

	schedule.Entries = make([]models.ScheduleEntry, 0, scheduleEntryEles.Length())

	var wg sync.WaitGroup
	var mutex sync.Mutex

	//Go through each class in the schedule
	scheduleEntryEles.Each(func(_ int, scheduleEntryEle *goquery.Selection) {
		wg.Add(1)

		go func() {
			defer wg.Done()
			scheduleEntry := parseScheduleEntry(scheduleEntryEle)
			mutex.Lock()
			defer mutex.Unlock()
			//Push the entry to the slice
			schedule.Entries = append(schedule.Entries, scheduleEntry)
		}()
	})

	wg.Wait()

	return schedule
}

// parseScheduleEntry takes in the HTML element for the schedule row, and parses it
// into a ScheduleEntry struct. It returns that plus a slice containing the marking
// periods this entry is active for.
func parseScheduleEntry(scheduleEntryEle *goquery.Selection) models.ScheduleEntry {
	//Make entry struct
	scheduleEntry := models.ScheduleEntry{}

	//Go through each td, using i to find what the text corresponds to
	scheduleEntryEle.Find("td").Each(func(i int, dataEle *goquery.Selection) {
		//Parse text, return if there is none
		text := strings.TrimSpace(dataEle.Text())
		if len(text) <= 0 {
			return
		}

		//Fill in data using i
		switch i {
		case 0:
			scheduleEntry.Class.Course = text
		case 1:
			scheduleEntry.Class.Name = text
		case 2:
			scheduleEntry.Class.Period = text
		case 3:
			scheduleEntry.Class.Teacher = text
		case 4:
			scheduleEntry.Class.Room = text
		case 5:
			scheduleEntry.Days = strings.Split(text, ", ")
		case 6:
			scheduleEntry.MarkingPeriods = strings.Split(text, ", ")
		case 7:
			scheduleEntry.Building = text
		case 8:
			scheduleEntry.Active = strings.EqualFold(text, "active")
		}
	})

	return scheduleEntry
}
