package stores

import "github.com/even44/JobsearchAPI/pkg/models"

type CompanyStore interface {
	AddCompany(company models.Company) (*models.Company, error)
	GetCompany(Id int) (*models.Company, error)
	ListCompanies() ([]models.Company, error)
	UpdateCompany(Id int, company models.Company) error
	RemoveCompany(id int) error
}