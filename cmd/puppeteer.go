package main

import (
	"encoding/csv"
	"fmt"
	"github.com/kenji-yamane/pq-deadlock/src/customerror"
	"github.com/kenji-yamane/pq-deadlock/src/puppeteer"
	"os"
	"strconv"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		customerror.CheckError(fmt.Errorf("not enough ports given as arguments"))
	}
	ports := os.Args[2:len(os.Args)]
	p := puppeteer.NewPuppeteer(ports)

	inputFile := os.Args[1]
	f, err := os.Open(inputFile)
	if err != nil {
		customerror.CheckError(fmt.Errorf("no file found"))
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		customerror.CheckError(fmt.Errorf("error parsing csv"))
	}
	for _, rec := range records {
		if len(rec) != 3 {
			customerror.CheckError(fmt.Errorf("there should be three columns only"))
		}
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			customerror.CheckError(fmt.Errorf("error parsing id"))
		}
		delay, err := strconv.Atoi(rec[2])
		if err != nil {
			customerror.CheckError(fmt.Errorf("error parsing delay"))
		}
		p.AddCommand(id, rec[1], time.Duration(delay)*time.Millisecond)
	}
	p.Execute()
}
