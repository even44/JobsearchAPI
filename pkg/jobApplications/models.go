package jobApplications

// Represents a job application
type JobApplication struct {
	Id int `json:"id" gorm:"unique;primaryKey;autoIncrement"`
	Position string `json:"position"`
	Company string `json:"company"`
	SearchDate string `json:"search_date"`
	Deadline string `json:"deadline"`
	Response bool `json:"response"`
	Interview bool `json:"interview"`
	Done bool `json:"done"`
	Link string `json:"link"`
	

}

