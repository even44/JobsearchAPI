package models


type Company struct {
	ID	 uint	   `json:"ID" gorm:"primarykey"`           
	UserID   uint      `json:"user_id"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Status   string    `json:"status"`
	Notes    string    `json:"notes"`
	Website  string    `json:"website"`
	Contacts []Contact `json:"contacts"`
}
