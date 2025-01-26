package models

import (
	"time"

	"gorm.io/gorm"
)

type JobApplication struct {
	gorm.Model
	Position   string    `json:"position"`
	CompanyId  int       `json:"company_id"`
	Company    Company   `json:"company"`
	SearchDate time.Time `json:"search_date"`
	Deadline   time.Time `json:"deadline"`
	Response   bool      `json:"response"`
	Interview  bool      `json:"interview"`
	Done       bool      `json:"done"`
	Link       string    `json:"link"`
}
