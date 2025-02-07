package stores

import (
	"fmt"

	"github.com/even44/JobsearchAPI/pkg/models"
)

func (s MariaDBStore) AddContact(contact models.Contact) (*models.Contact, error) {

	company, err := s.GetCompany(contact.CompanyID)
	if err != nil {
		s.logger.Printf("[WARN][ADD] Company does not exist")
		return nil, err
	}

	if company.UserID != contact.UserID {
		s.logger.Printf("[WARN][ADD] Company exists but does not belong to this user, aborting...")
		return nil, fmt.Errorf("invalid company id")
	}

	s.logger.Printf("[ADD] Contact with name %s will be created", contact.Name)
	result := s.db.Create(&contact)
	if result.Error != nil {
		return nil, result.Error
	}

	resConstact, err := s.GetContact(contact.ID)
	if err != nil {
		return nil, result.Error
	}
	s.logger.Printf("[ADD] Created contact with name %s", contact.Name)
	return resConstact, nil
}
func (s MariaDBStore) GetContact(id uint) (*models.Contact, error) {
	var contact models.Contact
	result := s.db.First(&contact, id)
	if result.Error != nil {
		return nil, result.Error
	}
	s.logger.Printf("[GET] Got contact with Name '%s' and id %d ", contact.Name, contact.ID)
	return &contact, nil
}
func (s MariaDBStore) ListContacts(userID uint) ([]models.Contact, error) {
	var contacts []models.Contact
	s.db.Preload("Company").Find(&contacts, &models.Contact{UserID: userID})
	s.logger.Printf("[LIST] Got %d companies", len(contacts))
	return contacts, nil
}
func (s MariaDBStore) UpdateContact(id uint, contact models.Contact) error {

	existingContact, err := s.GetContact(id)
	if err != nil {
		return err
	}

	if contact.CompanyID == 0 {
		contact.CompanyID = existingContact.CompanyID
	}

	company, err := s.GetCompany(contact.CompanyID)
	if err != nil {
		return err
	}
	if company.UserID != contact.UserID {
		s.logger.Printf("[WARN][ADD] Company exists but does not belong to this user, aborting...")
		return fmt.Errorf("invalid company id")
	}

	result := s.db.Save(&contact)
	if result.Error != nil {
		s.logger.Printf("[ERROR] %s", result.Error)
		return result.Error
	}
	s.logger.Printf(
		"[UPDATE] Updated contact with name '%s' and company_id '%d' with id %d ",
		contact.Name, contact.CompanyID, contact.ID)
	return nil
}
func (s MariaDBStore) RemoveContact(id uint) error {
	contact, err := s.GetContact(id)
	if err != nil {
		return err
	}
	result := s.db.Delete(contact)
	if result.Error != nil {
		return result.Error
	}
	s.logger.Printf(
		"[DELETE] Deleted contact with name '%s' and company_id '%d' with id %d ",
		contact.Name, contact.CompanyID, contact.ID)
	return nil
}
