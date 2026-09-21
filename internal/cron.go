package internal

import (
	"log"
	"time"
)

// Starts the daily cron job, it is triggered at 3am every day at 3am
func StartDailyCron(fn func()) {
	go func() {
		loc, err := time.LoadLocation("Europe/Oslo")
		if err != nil {
			log.Println("cron: falling back to UTC, couldn't load Europe/Oslo:", err.Error())
			loc = time.UTC
		}

		for {
			now := time.Now().In(loc)
			next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, loc)
			if !next.After(now) {
				next = next.Add(24 * time.Hour)
			}

			log.Printf("cron: next run scheduled for %s", next.Format(time.RFC3339))

			timer := time.NewTimer(next.Sub(now))
			<-timer.C

			log.Println("cron: run starting")
			start := time.Now()

			fn()

			log.Printf("cron: run finished in %s", time.Since(start))
		}
	}()
}
