package stores

import (
	"github.com/even44/JobsearchAPI/pkg/models"
)

func (s MariaDBStore) AddApplication(application models.JobApplication) (*models.JobApplication, error) {

	s.logger.Println("[ADD] Creating Application")
	result := s.db.Omit("Company").Create(&application)
	if result.Error != nil {
		return nil, result.Error
	}

	s.logger.Println("[ADD] Getting created application")
	resApplication, err := s.GetApplication(application.ID)
	if err != nil {
		return nil, err
	}

	s.logger.Printf("[ADD] Created job application with position '%s' at '%s' with id %d ", resApplication.Position, resApplication.Company.Name, resApplication.ID)
	return resApplication, nil

}
func (s MariaDBStore) GetApplication(Id uint) (*models.JobApplication, error) {
	var application models.JobApplication
	result := s.db.Preload("Company.Contacts").Preload("Company").First(&application, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	s.logger.Printf("[GET] Got job application with position '%s' and company '%s' with id %d ", application.Position, application.Company.Name, application.ID)
	return &application, nil
}
func (s MariaDBStore) ListApplications() ([]models.JobApplication, error) {
	var applications []models.JobApplication
	s.db.Preload("Company.Contacts").Preload("Company").Find(&applications)
	s.logger.Printf("[LIST] Got %d job applications", len(applications))
	return applications, nil
}
func (s MariaDBStore) UpdateApplication(id uint, application models.JobApplication) error {
	result := s.db.Omit("Company").Save(&application)
	if result.Error != nil {
		return result.Error
	}
	s.logger.Printf("[UPDATE] Updated job application with position '%s' with id %d ", application.Position, application.ID)
	return nil
}
func (s MariaDBStore) RemoveApplication(id uint) error {
	application, err := s.GetApplication(id)
	if err != nil {
		return err
	}
	result := s.db.Delete(application)
	if result.Error != nil {
		return result.Error
	}
	s.logger.Printf("[DELETE] Deleted job application with position '%s' and company '%s' with id %d ", application.Position, application.Company.Name, application.ID)
	return nil
}
