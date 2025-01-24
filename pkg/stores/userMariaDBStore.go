package stores

import "github.com/even44/JobsearchAPI/pkg/models"

func (s MariaDBStore) Add(user *models.User) error {
	err := s.db.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (s MariaDBStore) Get(email string) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, models.User{Email: email}).Error
	if err != nil {
		return nil, err
	}
	return &user, nil

}
