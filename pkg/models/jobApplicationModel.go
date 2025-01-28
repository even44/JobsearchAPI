package models

type JobApplication struct {
	ID         uint      `gorm:"primarykey"`
	UserID     uint      `json:"user_id" gorm:"index:idx_application"`
	CompanyID  uint      `json:"company_id" gorm:"index:idx_application"`
	Position   string    `json:"position"`
	Company    Company   `json:"company"`
	SearchDate CivilTime `json:"search_date" gorm:"type:text"`
	Deadline   CivilTime `json:"deadline" gorm:"type:text"`
	Response   bool      `json:"response"`
	Interview  bool      `json:"interview"`
	Done       bool      `json:"done"`
	Link       string    `json:"link"`
	ContactID  uint      `json:"contact_id"`
	Contact    Contact   `json:"contact"`
}
