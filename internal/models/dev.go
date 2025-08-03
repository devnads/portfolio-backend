package models

import "time"

type Dev struct {
	ID string `json:"id"`
	Username string `json:"username"`
	ProfileImage string `json:"profile_image"`
	Roles []string `json:"roles"`
	Address string `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	Discord string `json:"discord"`
	Twitter string `json:"twitter"`
	Github string `json:"github"`
}