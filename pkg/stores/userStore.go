package stores

import "github.com/even44/JobsearchAPI/pkg/models"

type UserStore interface {
	Add(user *models.User) error
	Get(email string) (*models.User, error)
	//Update(*models.User) error
	//Remove(id int) error
}
