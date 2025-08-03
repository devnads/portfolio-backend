package models

import "time"

type Project struct {
	ID          string    `json:"id"`
	DevID	    string    `json:"dev_id"`
	MissionID   string    `json:"mission_id"`
	Name        string    `json:"name"`
	Image       string    `json:"image"`
	Categories  []string  `json:"categories"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}