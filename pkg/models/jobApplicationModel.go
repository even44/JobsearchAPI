package models

type JobApplication struct {
	Id         int     `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Position   string  `json:"position"`
	CompanyId  int     `json:"company_id"`
	Company    Company `json:"company"`
	SearchDate string  `json:"search_date"`
	Deadline   string  `json:"deadline"`
	Response   bool    `json:"response"`
	Interview  bool    `json:"interview"`
	Done       bool    `json:"done"`
	Link       string  `json:"link"`
}
