package stores

import (
	"errors"

	"github.com/even44/JobsearchAPI/pkg/models"
	"gorm.io/gorm"
)

func (s MariaDBStore) AddContact(contact models.Contact) (*models.Contact, error) {

	err := s.db.First(&models.Contact{}, models.Contact{Name: contact.Name, UserID: contact.UserID}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Printf("[WARN][ADD] Contact with name %s already exists and will not be created", contact.Name)
		return nil, err
	}
	s.logger.Printf("[ADD] Contact with name %s does not exist and will be created", contact.Name)
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
	s.db.Find(&contacts, &models.Contact{UserID: userID})
	s.logger.Printf("[LIST] Got %d companies", len(contacts))
	return contacts, nil
}
func (s MariaDBStore) UpdateContact(id uint, contact models.Contact) error {

	var existingContact *models.Contact
	err := s.db.First(&existingContact, models.Contact{Name: contact.Name, UserID: contact.UserID}).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		if existingContact.ID != id {
			s.logger.Printf(
				"[WARN][UPDATE] A different contact with name %s already exists and will not be updated",
				contact.Name)
			return err
		}
	}

	if contact.CompanyID == 0 {
		contact.CompanyID = existingContact.CompanyID
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
