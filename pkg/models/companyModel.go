package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	UserID   uint
	Name     string    `json:"name" gorm:"index:idx_name,unique"`
	Location string    `json:"location" gorm:"index:idx_name,unique"`
	Status   string    `json:"status"`
	Notes    string    `json:"notes"`
	Website  string    `json:"website"`
	Contacts []Contact `json:"contacts"`
}
