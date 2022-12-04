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
)

type Journal struct {
	indexPath, jsonPath string
	Name                string
	MileageTotal        int      `json:"mileage_total"`
	ExpenseTotal        float64  `json:"expense_total"`
	Entries             []*Entry `json:"entries"`
	States              []string `json:"all_states,omitempty"`
	USParks             []string `json:"all_us_parks,omitempty"`
	Provinces           []string `json:"all_provinces,omitempty"`
	CAParks             []string `json:"all_ca_parks,omitempty"`
}

func New(name string) (j *Journal, err error) {
	j = &Journal{
		indexPath: fmt.Sprintf("%s/README.md", name),
		jsonPath:  fmt.Sprintf("%s/journal.json", name),
		Name:      name,
	}

	if err = j.unmarshal(j.jsonPath); err != nil {
		return nil, err
	}

	for _, e := range j.Entries {
		pEnd := 60.0
		if p, err := j.previousEntry(e); err == nil {
			pEnd = p.BudgetEnd
		}

		e.updateNewEntry(pEnd)
		j.updateJournal(e)

		e.allLocations = [][]string{j.States, j.USParks, j.Provinces, j.CAParks}
	}

	return j, nil
}

func (j *Journal) updateJournal(e *Entry) {
	j.MileageTotal += e.Mileage
	j.ExpenseTotal += e.DailyExpenses

	j.updateLocationList(e, states)
	j.updateLocationList(e, usParks)
	j.updateLocationList(e, provinces)
	j.updateLocationList(e, caParks)

}

func (j *Journal) updateLocationList(e *Entry, category listCategory) {
	var eList, jList *[]string
	switch category {
	case states:
		eList = &e.States
		jList = &j.States
	case usParks:
		eList = &e.USParks
		jList = &j.USParks
	case provinces:
		eList = &e.Provinces
		jList = &j.Provinces
	case caParks:
		eList = &e.CAParks
		jList = &j.CAParks
	}

	if eList != nil {
		for _, item := range *eList {
			var found bool
			for _, exItem := range *jList {
				if item == exItem {
					found = true
				}
			}
			if !found {
				*jList = append(*jList, item)
			}
		}
	}
}

func (j *Journal) unmarshal(jsonPath string) (err error) {
	var jsonFile *os.File
	if jsonFile, err = os.Open(jsonPath); err != nil {
		return
	}
	defer func() { _ = jsonFile.Close() }()

	var bytes []byte
	if bytes, err = io.ReadAll(jsonFile); err != nil {
		return
	}

	if err = json.Unmarshal(bytes, &j); err != nil {
		return
	}

	// zero values are needed for totals and tallies after unmarshalling
	j.MileageTotal = 0
	j.ExpenseTotal = 0
	j.States = make([]string, 0)
	j.USParks = make([]string, 0)
	j.Provinces = make([]string, 0)
	j.CAParks = make([]string, 0)

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

func (j *Journal) WriteIndex(e *Entry) error {
	file, err := os.OpenFile(fmt.Sprintf("%s/README.md", j.Name), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
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

func (j *Journal) HasFile(e *Entry, f fileType) bool {
	name := fmt.Sprintf(f.format(), j.Name, e.Name)
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}
