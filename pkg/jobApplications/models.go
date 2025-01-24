package jobApplications

// Represents a job application
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

type Company struct {
	Id       int       `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Status   string    `json:"status"`
	Notes    string    `json:"notes"`
	Website  string    `json:"website"`
	Contacts []Contact `json:"contacts"`
}

type Contact struct {
	Id        int    `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Name      string `json:"name"`
	CompanyId int    `json:"company_id"`
	Email     string `json:"email"`
	Phone     int    `json:"phone"`
}
