package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	UserID    uint
	Name      string `json:"name"`
	CompanyID uint   `json:"company_id"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
}
