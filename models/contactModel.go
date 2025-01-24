package models

type Contact struct {
	Id        int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name      string `json:"name"`
	CompanyId int    `json:"company_id"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
}
