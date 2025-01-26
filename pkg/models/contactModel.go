package models

import "gorm.io/gorm"

type Contact struct {
	gorm.Model
	UserID    uint   `json:"user_id" gorm:"index:idx_companycontact"`
	CompanyID uint   `json:"company_id" gorm:"index:idx_companycontact"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
	
}
