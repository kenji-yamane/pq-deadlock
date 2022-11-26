package main

import (
	"github.com/kenji-yamane/pq-deadlock/src"
	"github.com/kenji-yamane/pq-deadlock/src/network"
)

func main() {
	serverCh := make(chan string)
	go network.Serve(serverCh, src.SharedResourcePort)
	for {
		_, valid := <-serverCh
		if !valid {
			continue
		}
	}
}
