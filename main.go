package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	terminate := make(chan bool)
	polled := make(chan bool)
	pollingWorker(terminate, polled)
	pruningWorker(terminate, polled)

	// create a channel to receive incoming OS interrupts (such as Ctrl-C):
	osInterruptChannel := make(chan os.Signal, 1)
	signal.Notify(osInterruptChannel, os.Interrupt)

	// block execution until an OS signal (such as Ctrl-C) is received:
	<-osInterruptChannel
	terminate <- true
}

func pollingWorker(terminate <-chan bool, polled chan<- bool) {
	go func() {
		for {
			select {
			case <-terminate:
				fmt.Println("Quitting")
				return
			default:
				fmt.Println("Polling")
				time.Sleep(2 * time.Second)
				polled <- true
			}
		}
	}()
}

func pruningWorker(terminate <-chan bool, polled <-chan bool) {
	go func() {
		for {
			select {
			case <-terminate:
				fmt.Println("Quitting")
				return
			case <-polled:
				fmt.Println("Pruning")
				time.Sleep(2 * time.Second)
			}
		}
	}()
}
