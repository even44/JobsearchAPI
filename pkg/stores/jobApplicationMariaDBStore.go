package stores

import (
	"fmt"

	"github.com/even44/JobsearchAPI/pkg/models"
)

func (s MariaDBStore) AddApplication(application models.JobApplication) (*models.JobApplication, error) {

	company, err := s.GetCompany(application.CompanyID)
	if err != nil {
		s.logger.Printf("[WARN][ADD] Company does not exist")
		return nil, err
	}

	contact, err := s.GetContact(application.ContactID)
	if err != nil {
		s.logger.Printf("[WARN][ADD] Contact does not exist")
	}

	if company.UserID != application.UserID {
		s.logger.Printf("[WARN][ADD] Company exists but does not belong to this user, aborting...")
		return nil, fmt.Errorf("invalid company id")
	}

	if contact.UserID != application.UserID {
		s.logger.Printf("[WARN][ADD] Contact exists but does not belong to this user, aborting...")
		return nil, fmt.Errorf("invalid contact id")
	}

	if contact.CompanyID != application.CompanyID {
		s.logger.Printf("[WARN][ADD] Contact exists but does not belong to this company, aborting...")
		return nil, fmt.Errorf("invalid contact id")
	}

	s.logger.Println("[ADD] Creating Application")
	result := s.db.Omit("Company", "Contact").Create(&application)
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
	result := s.db.Preload("Contact").Preload("Company.Contacts").Preload("Company").First(&application, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	s.logger.Printf("[GET] Got job application with position '%s' and company '%s' with id %d ", application.Position, application.Company.Name, application.ID)
	return &application, nil
}
func (s MariaDBStore) ListApplications(userID uint) ([]models.JobApplication, error) {
	var applications []models.JobApplication
	s.db.Preload("Contact").Preload("Company.Contacts").Preload("Company").Find(&applications, &models.JobApplication{UserID: userID})
	s.logger.Printf("[LIST] Got %d job applications", len(applications))
	return applications, nil
}
func (s MariaDBStore) UpdateApplication(id uint, application models.JobApplication) error {

	existingApplication, err := s.GetApplication(id)
	if err != nil {
		return err
	}

	if application.CompanyID == 0 {
		application.CompanyID = existingApplication.CompanyID
	}

	company, err := s.GetCompany(application.CompanyID)
	if err != nil {
		return err
	}

	contact, err := s.GetContact(application.ContactID)
	if err != nil {
		s.logger.Printf("[WARN][UPDATE] Contact does not exist")
	}

	if company.UserID != application.UserID {
		s.logger.Printf("[WARN][UPDATE] Company exists but does not belong to this user, aborting...")
		return fmt.Errorf("invalid company id")
	}

	if contact.UserID != application.UserID {
		s.logger.Printf("[WARN][UPDATE] Contact exists but does not belong to this user, aborting...")
		return fmt.Errorf("invalid contact id")
	}

	if contact.CompanyID != application.CompanyID {
		s.logger.Printf("[WARN][UPDATE] Contact exists but does not belong to this company, aborting...")
		return fmt.Errorf("invalid contact id")
	}

	result := s.db.Omit("Company", "Contact").Save(&application)
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
