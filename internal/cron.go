package internal

import (
	"fmt"
	"time"
)

// Starts the daily cron job, it is triggered at 3am every day at 3am
func StartDailyCron(fn func()) {
	go func() {
		loc, err := time.LoadLocation("Europe/Oslo")
		if err != nil {
			fmt.Println("cron: falling back to UTC, couldn't load Europe/Oslo:", err.Error())
			loc = time.UTC
		}

		for {
			now := time.Now().In(loc)
			next := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, loc)
			if !next.After(now) {
				next = next.Add(24 * time.Hour)
			}

			timer := time.NewTimer(next.Sub(now))
			<-timer.C

			fn()
		}
	}()
}
