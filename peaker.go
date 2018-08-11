package main

import (
	"fmt"
	"log"
	"lukweb.de/mailPeaker/peaker"
	_ "net/smtp"
	"time"
)

var config *peaker.Config

func main() {
	localConfig, err := peaker.ReadConfig()
	if err != nil {
		panic(err)
	}
	config = localConfig

	if config.Dev {
		log.Println("Dev mode is enabled => instant run")
		tickWithInterval()
	} else {
		waitForNextInterval()
	}
}

func waitForNextInterval() {
	nextTime := time.Now().Truncate(config.Interval).Add(config.Interval)

	log.Printf("First run at %v\n", nextTime)
	time.Sleep(time.Until(nextTime))

	tickWithInterval()
}

func tickWithInterval() {
	exec(time.Now())
	ticker := time.NewTicker(config.Interval)
	for t := range ticker.C {
		if !exec(t) {
			break
		}
	}
}

func exec(t time.Time) bool {
	fmt.Println("test : -- " + t.String())
	return true
}
