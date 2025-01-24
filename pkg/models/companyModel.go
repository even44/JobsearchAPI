package models

type Company struct {
	Id       int       `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name     string    `json:"name" gorm:"index:idx_name,unique"`
	Location string    `json:"location" gorm:"index:idx_name,unique"`
	Status   string    `json:"status"`
	Notes    string    `json:"notes"`
	Website  string    `json:"website"`
	Contacts []Contact `json:"contacts"`
}