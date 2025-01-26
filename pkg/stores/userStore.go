package stores

import "github.com/even44/JobsearchAPI/pkg/models"

type UserStore interface {
	AddUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserById(id int) (*models.User, error)
	//Update(*models.User) error
	//Remove(id int) error
}
