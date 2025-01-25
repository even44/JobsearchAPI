package stores

import "github.com/even44/JobsearchAPI/pkg/models"

type UserStore interface {
	Add(user *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetById(id int) (*models.User, error)
	//Update(*models.User) error
	//Remove(id int) error
}
