package stores

import "github.com/even44/JobsearchAPI/pkg/models"

type CompanyStore interface {
	AddCompany(company models.Company) (*models.Company, error)
	GetCompany(Id uint) (*models.Company, error)
	ListCompanies() ([]models.Company, error)
	UpdateCompany(Id uint, company models.Company) error
	RemoveCompany(id uint) error
}