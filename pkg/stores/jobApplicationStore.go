package stores

import "github.com/even44/JobsearchAPI/pkg/models"

type JobApplicationStore interface {
	AddApplication(jobApplication models.JobApplication) (*models.JobApplication, error)
	GetApplication(id int) (*models.JobApplication, error)
	ListApplications() ([]models.JobApplication, error)
	UpdateApplication(id int, jobApplication models.JobApplication) error
	RemoveApplication(id int) error
}