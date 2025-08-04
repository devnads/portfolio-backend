package models

import "time"

type Mission struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Round       int       `json:"round"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Projects    []Project `json:"projects"`
}
