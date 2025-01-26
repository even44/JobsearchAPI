package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	Name      string `json:"name"`
	CompanyId int    `json:"company_id"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
}
