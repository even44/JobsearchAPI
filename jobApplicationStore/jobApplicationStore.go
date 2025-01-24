package jobapplicationstore

import "github.com/even44/JobsearchAPI/models"

type JobApplicationStore interface {
	AddApplication(id int, jobApplication models.JobApplication) (*models.JobApplication, error)
	GetApplication(id int) (*models.JobApplication, error)
	ListApplications() ([]models.JobApplication, error)
	UpdateApplication(id int, jobApplication models.JobApplication) error
	RemoveApplication(id int) error
	AddCompany(id int, company models.Company) (*models.Company, error)
	GetCompany(Id int) (*models.Company, error)
	ListCompanies() ([]models.Company, error)
	UpdateCompany(Id int, company models.Company) error
	RemoveCompany(id int) error
	AddContact(company_id int, contact models.Contact) (*models.Contact, error)
	GetContact(id int) (*models.Contact, error)
	ListContacts() ([]models.Contact, error)
	UpdateContact(id int, contact models.Contact) error
	RemoveContact(id int) error
}