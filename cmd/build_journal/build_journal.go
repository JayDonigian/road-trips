package main

import (
	"flag"
	"github.com/JayDonigian/road-trips/journal"
	"log"
)

func main() {
	nameArg := flag.String("name", "", "the name of the journal")
	flag.Parse()
	journalName := *nameArg

	if *nameArg == "" {
		log.Fatalf("This script requires the journal name as a parameter.\n")
	}

	j, err := journal.New(journalName)
	if err != nil {
		log.Fatalf("ERROR: while creating journal - %s", err.Error())
	}

	for _, e := range j.MissingEntries() {
		err = e.WriteFile(j)
		if err != nil {
			log.Fatalf("ERROR: while creating from template file - %s", err)
		}

		err = j.WriteIndex(e)
		if err != nil {
			log.Fatalf("ERROR: while creating from template file - %s", err)
		}
	}

	err = j.Save()
	if err != nil {
		log.Printf("WARNING: while saving journal - %s", err.Error())
	}
}
