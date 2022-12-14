package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/JayDonigian/road-trips/journal"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	if e, err := entryFromFile(); err == nil {
		var bytes []byte
		if bytes, err = json.MarshalIndent(e, "", "    "); err == nil {
			_ = os.WriteFile(fmt.Sprintf("2016/%s.json", e.Name), bytes, 0644)
		}
	}
}

func entryFromFile() (*journal.Entry, error) {
	var file *os.File
	var err error
	if file, err = os.Open("2016/completed_form.md"); err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	info := &journal.Entry{}

	var expensing, locatingStates, locatingProvinces, locatingUSParks, locatingCAParks, journaling, exit bool
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if exit {
			break
		}
		line := scanner.Text()
		switch {
		case strings.Contains(line, "Date (mm-dd)"):
			info.Name = strings.Split(line, "`")[1]
			var t time.Time
			if t, err = time.Parse("01-02", info.Name); err == nil {
				info.Date = time.Date(2016, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
			}
		case strings.Contains(line, "Mileage"):
			info.Mileage, _ = strconv.Atoi(strings.Split(line, "`")[1])
		case strings.Contains(line, "Start Short Name"):
			info.Start.Short = strings.Split(line, "`")[1]
		case strings.Contains(line, "Start Long Name"):
			info.Start.Long = strings.Split(line, "`")[1]
		case strings.Contains(line, "Start Emoji"):
			info.Start.Emoji = strings.Split(line, "`")[1]
		case strings.Contains(line, "End Short Name"):
			info.End.Short = strings.Split(line, "`")[1]
		case strings.Contains(line, "End Long Name"):
			info.End.Long = strings.Split(line, "`")[1]
		case strings.Contains(line, "End Emoji"):
			info.End.Emoji = strings.Split(line, "`")[1]
		case strings.Contains(line, "### Expenses"):
			expensing = true
			info.Expenses = []journal.Expense{}
		case expensing:
			if strings.Contains(line, "`") {
				cost, _ := strconv.ParseFloat(strings.Split(line, "`")[3], 64)
				info.Expenses = append(info.Expenses, journal.Expense{Item: strings.Split(line, "`")[1], Cost: cost})
			}
			if line == "" {
				expensing = false
			}
		case strings.Contains(line, "* States:"):
			locatingStates = true
			info.States = make([]string, 0)
		case locatingStates:
			if strings.Contains(line, "`") {
				info.States = append(info.States, strings.Split(line, "`")[1])
			}
			if strings.Contains(line, "* National Parks:") {
				locatingStates = false
				locatingUSParks = true
				info.USParks = make([]string, 0)
			}
			if line == "" {
				locatingStates = false
			}
		case locatingUSParks:
			if strings.Contains(line, "`") {
				info.USParks = append(info.USParks, strings.Split(line, "`")[1])
			}
			if strings.Contains(line, "* Provinces:") {
				locatingUSParks = false
				locatingProvinces = true
				info.Provinces = make([]string, 0)
			}
			if line == "" {
				locatingUSParks = false
			}
		case locatingProvinces:
			if strings.Contains(line, "`") {
				info.Provinces = append(info.Provinces, strings.Split(line, "`")[1])
			}
			if strings.Contains(line, "* Canadian National Parks:") {
				locatingProvinces = false
				locatingCAParks = true
				info.CAParks = make([]string, 0)
			}
			if line == "" {
				locatingProvinces = false
			}
		case locatingCAParks:
			if strings.Contains(line, "`") {
				info.Provinces = append(info.Provinces, strings.Split(line, "`")[1])
			}
			if line == "" {
				locatingCAParks = false
			}
		case strings.Contains(line, "### Emoji Story"):
			info.EmojiStory = strings.Split(line, "`")[1]
		case strings.Contains(line, "### Journal Entry"):
			journaling = true
			info.JournalEntry = make([]string, 0)
		case journaling:
			if line == "" {
				journaling = false
				exit = true
				break
			}
			if !strings.Contains(line, "```") {
				info.JournalEntry = append(info.JournalEntry, line)
			}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	return info, nil
}
