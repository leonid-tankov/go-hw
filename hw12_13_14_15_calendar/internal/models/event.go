package models

import "time"

type Event struct {
	ID             string
	Title          string
	Timestamp      time.Time
	Duration       time.Duration
	Description    string
	CustomerID     string
	NotifyDuration time.Duration
}
