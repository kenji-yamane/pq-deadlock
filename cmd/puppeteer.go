package main

import (
	"fmt"
	"github.com/kenji-yamane/pq-deadlock/src/customerror"
	"github.com/kenji-yamane/pq-deadlock/src/puppeteer"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		customerror.CheckError(fmt.Errorf("not enough ports given as arguments"))
	}
	ports := os.Args[1:len(os.Args)]

	p := puppeteer.NewPuppeteer(ports)
	p.AddCommand(1, "ask 1 2", 3*time.Second)
	p.AddCommand(2, "ask 1 3", 3*time.Second)
	p.AddCommand(3, "ask 1 1", 3*time.Second)
	p.AddCommand(1, "detect", 3*time.Second)
	p.Execute()
}
