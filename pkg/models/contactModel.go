package models


type Contact struct {
	ID 	  uint    `json:"ID" gorm:"primarykey"`
	UserID    uint    `json:"user_id" gorm:"index:idx_companycontact"`
	CompanyID uint    `json:"company_id" gorm:"index:idx_companycontact"`
	Company   Company `json:"company"`
	Name      string  `json:"name"`
	Email     string  `json:"email"`
	Phone     int     `json:"phone"`
	
}
