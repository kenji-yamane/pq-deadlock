package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kenji-yamane/pq-deadlock/src"
	"github.com/kenji-yamane/pq-deadlock/src/clock"
	"github.com/kenji-yamane/pq-deadlock/src/customerror"
	"github.com/kenji-yamane/pq-deadlock/src/process"
)

func main() {
	if len(os.Args) < 4 {
		customerror.CheckError(fmt.Errorf("not enough ports given as arguments"))
	}
	myId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		customerror.CheckError(fmt.Errorf("first argument should be a number representing the sequential process ID"))
	}
	ports := os.Args[2:len(os.Args)]

	p := process.NewProcess(
		myId,
		ports,
		clock.NewScalarClock(myId),
	)

	terminalCh := make(chan string)
	go src.ReadInput(terminalCh)

	serverCh := make(chan string)
	go p.Serve(serverCh)

	for {
		select {
		case command, valid := <-terminalCh:
			if !valid {
				break
			}
			p.InterpretCommand(command)
		case msg, valid := <-serverCh:
			if !valid {
				break
			}
			p.InterpretMessage(msg)
		default:
		}
		time.Sleep(time.Second * 1)
	}
}
