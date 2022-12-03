package journal

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type Journal struct {
	indexPath    string
	Name         string
	MileageTotal int      `json:"mileage_total"`
	ExpenseTotal float64  `json:"expense_total"`
	Entries      []*Entry `json:"entries"`
}

func New(name string) (*Journal, error) {
	j := &Journal{Name: name}
	err := j.unmarshal(fmt.Sprintf("%s/journal.json", j.Name))
	if err != nil {
		return nil, err
	}

	j.indexPath = fmt.Sprintf("%s/index.md", j.Name)

	var t time.Time
	for _, e := range j.Entries {
		t, err = time.Parse("01-02", e.Name)
		if err != nil {
			return nil, err
		}
		e.Date = time.Date(2016, t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())

		if e.DailyExpenses == 0 {
			for _, expense := range e.Expenses {
				e.DailyExpenses += expense.Cost
			}
		}

		var p *Entry
		var pEnd float64
		if p, err = j.previousEntry(e); err == nil {
			pEnd = p.BudgetEnd
		}

		e.BudgetStart = pEnd + 60.00
		e.BudgetEnd = e.BudgetStart - e.DailyExpenses

		j.MileageTotal += e.Mileage
		j.ExpenseTotal += e.DailyExpenses
	}

	return j, nil
}

func (j *Journal) unmarshal(jsonPath string) error {
	jsonFile, err := os.Open(jsonPath)
	if err != nil {
		return err
	}
	defer func() { _ = jsonFile.Close() }()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &j)
	if err != nil {
		return err
	}
	return nil
}

func (j *Journal) previousEntry(entry *Entry) (*Entry, error) {
	p := entry.Date.AddDate(0, 0, -1)
	for _, e := range j.Entries {
		if e.Date == p {
			return e, nil
		}
	}
	return nil, errors.New("unable to find a previous entry")
}

func (j *Journal) MissingEntries() []*Entry {
	var missing []*Entry
	for _, e := range j.Entries {
		if !j.HasFile(e, dayMap) {
			log.Printf("WARNING: day map for %s does not exist\n", e.Name)
		}
		if !j.HasFile(e, totalMap) {
			log.Printf("WARNING: total map for %s does not exist\n", e.Name)
		}
		if !j.HasFile(e, entry) {
			missing = append(missing, e)
		}
	}
	return missing
}

func (j *Journal) Write(e *Entry) error {
	destination, err := os.Create(fmt.Sprintf(fileType(entry).format(), j.Name, e.Name))
	if err != nil {
		return err
	}
	defer func() { _ = destination.Close() }()

	writer := bufio.NewWriter(destination)
	defer func() { _ = writer.Flush() }()

	lines := e.Write()
	for _, line := range j.TotalTripStats(e) {
		lines = append(lines, line)
	}
	for _, line := range e.PrevNextLinks() {
		lines = append(lines, line)
	}
	for _, line := range lines {
		_, _ = writer.WriteString(line + "\n")
	}

	return nil
}

func (j *Journal) WriteIndex(e *Entry) error {
	file, err := os.OpenFile(fmt.Sprintf("%s/index.md", j.Name), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), e.Name) {
			return nil
		}
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	_, err = file.WriteString(fmt.Sprintf("%s\n", e.Index()))
	if err != nil {
		return err
	}

	return nil
}

func (j *Journal) Save() error {
	jsonString, _ := json.MarshalIndent(j, "", "    ")
	err := os.WriteFile(fmt.Sprintf("%s/journal.json", j.Name), jsonString, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (j *Journal) TotalTripStats(e *Entry) []string {
	return []string{
		"## Trip Statistics\n",
		fmt.Sprintf("* **Total Distance:** %d miles", j.MileageTotal),
		fmt.Sprintf("* **Total Budget Spent:** $%.2f", j.ExpenseTotal),
		"* **U.S. States**",
		"  * New Hampshire",
		"  * Maine",
		"* **Canadian Provinces**",
		"  * Nova Scotia",
		"* **National Parks**",
		"  * Acadia\n",
		fmt.Sprintf("![total trip from Fremont to %s](%s \"total trip map\")\n", e.End.Short, e.RelativePathToFile(totalMap)),
	}
}

type fileType int

const (
	entry = iota
	dayMap
	bikeMap
	totalMap
)

func (ft fileType) format() string {
	switch ft {
	case entry:
		return "%s/entries/%s.md"
	case dayMap:
		return "%s/maps/day/%s.png"
	case bikeMap:
		return "%s/maps/bike/%s.png"
	case totalMap:
		return "%s/maps/total/%s-total.png"
	default:
		return ""
	}
}

func (ft fileType) formatPathRelativeToEntry() string {
	switch ft {
	case entry:
		return "%s.md"
	case dayMap:
		return "../maps/day/%s.png"
	case bikeMap:
		return "../maps/bike/%s.png"
	case totalMap:
		return "../maps/total/%s-total.png"
	default:
		return ""
	}
}

func (j *Journal) HasFile(e *Entry, f fileType) bool {
	name := fmt.Sprintf(f.format(), j.Name, e.Name)
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}
