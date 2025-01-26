package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	UserID   uint      `json:"user_id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Status   string    `json:"status"`
	Notes    string    `json:"notes"`
	Website  string    `json:"website"`
	Contacts []Contact `json:"contacts"`
}
