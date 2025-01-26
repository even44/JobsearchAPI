package stores

import "github.com/even44/JobsearchAPI/pkg/models"

type ContactStore interface {
	AddContact(contact models.Contact) (*models.Contact, error)
	GetContact(id int) (*models.Contact, error)
	ListContacts() ([]models.Contact, error)
	UpdateContact(id int, contact models.Contact) error
	RemoveContact(id int) error
}
