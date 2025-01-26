package stores

import "github.com/even44/JobsearchAPI/pkg/models"

type ContactStore interface {
	AddContact(contact models.Contact) (*models.Contact, error)
	GetContact(id uint) (*models.Contact, error)
	ListContacts(userID uint) ([]models.Contact, error)
	UpdateContact(id uint, contact models.Contact) error
	RemoveContact(id uint) error
}
