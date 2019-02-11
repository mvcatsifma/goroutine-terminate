package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	terminate := make(chan bool)
	polled := make(chan time.Time, 1)
	pollingWorker(terminate, polled)
	pruningWorker(terminate, polled)

	// create a channel to receive incoming OS interrupts (such as Ctrl-C):
	osInterruptChannel := make(chan os.Signal, 1)
	signal.Notify(osInterruptChannel, os.Interrupt)

	// block execution until an OS signal (such as Ctrl-C) is received:
	<-osInterruptChannel
	terminate <- true
}

func pollingWorker(terminate <-chan bool, polled chan<- time.Time) {
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-terminate:
				fmt.Println("Quitting")
				return
			case <-ticker.C:
				ts := time.Now().UTC()
				fmt.Println("Polling")
				select {
				case polled <- ts:
				default:
				}
			}
		}
	}()
}

func pruningWorker(terminate <-chan bool, polled <-chan time.Time) {
	go func() {
		for {
			select {
			case <-terminate:
				fmt.Println("Quitting")
				return
			case _ = <-polled:
				fmt.Println("Pruning")
				time.Sleep(2 * time.Second)
			}
		}
	}()
}
