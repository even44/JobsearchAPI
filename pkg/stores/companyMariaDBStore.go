package stores

import (
	"errors"

	"github.com/even44/JobsearchAPI/pkg/models"
	"gorm.io/gorm"
)

func (s MariaDBStore) AddCompany(company models.Company) (*models.Company, error) {

	err := s.db.First(&models.Company{}, models.Company{Name: company.Name, Location: company.Location}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Printf("[WARN][ADD] Company with name %s and location %s already exists and will not be created", company.Name, company.Location)
		s.db.First(&company, &models.Company{Name: company.Name, Location: company.Location})
		return &company, err
	}

	s.logger.Printf("[ADD] Company with name %s and location %s does not exist and will be created", company.Name, company.Location)
	result := s.db.Omit("Contacts").Create(&company)
	if result.Error != nil {
		return &company, result.Error
	}

	resCompany, err := s.GetCompany(company.ID)
	if err != nil {
		return nil, result.Error
	}
	s.logger.Printf("[ADD] Created company with Name '%s' and location %s ", company.Name, company.Location)
	return resCompany, nil
}
func (s MariaDBStore) GetCompany(Id uint) (*models.Company, error) {
	var company models.Company
	result := s.db.Preload("Contacts").First(&company, Id)
	if result.Error != nil {
		s.logger.Printf("[ERROR][GET] Could not find company with id %d", Id)
		return nil, result.Error
	}
	s.logger.Printf("[GET] Got company with Name '%s' and id %d ", company.Name, company.ID)
	return &company, nil
}
func (s MariaDBStore) ListCompanies() ([]models.Company, error) {
	var companies []models.Company
	s.db.Preload("Contacts").Find(&companies)
	s.logger.Printf("[LIST] Got %d companies", len(companies))
	return companies, nil
}
func (s MariaDBStore) UpdateCompany(Id uint, company models.Company) error {
	var existingCompany *models.Company
	err := s.db.First(&existingCompany, models.Company{Name: company.Name, Location: company.Location}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if existingCompany.ID != Id {
			s.logger.Printf("[WARN][UPDATE] Company with name %s and location %s already exists and will not be updated", company.Name, company.Location)
			return err
		}
	}

	result := s.db.Save(&company)
	if result.Error != nil {
		s.logger.Printf("[ERROR] %s", result.Error)
		return result.Error
	}
	s.logger.Printf("[UPDATE] Updated company with name '%s' and location '%s' with id %d ", company.Name, company.Location, company.ID)
	return nil
}
func (s MariaDBStore) RemoveCompany(id uint) error {
	company, err := s.GetCompany(id)
	if err != nil {
		return err
	}
	result := s.db.Delete(company)
	if result.Error != nil {
		return result.Error
	}
	s.logger.Printf("[DELETE] Deleted company with name '%s' and location '%s' with id %d ", company.Name, company.Location, company.ID)
	return nil
}
