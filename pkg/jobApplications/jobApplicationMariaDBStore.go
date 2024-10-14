package jobApplications

import (
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "log"
)

type MariaDBStore struct {
	db *gorm.DB
}

func NewMariaDBStore() *MariaDBStore{
	db, err := connectToMariaDB()
	if err != nil {
		return nil
	}
	return &MariaDBStore{
		db,
	}
}

func connectToMariaDB() (*gorm.DB, error) {
    dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    return db, nil
}


func (s MariaDBStore) Add(id int, application JobApplication) error {
	result := s.db.Create(application)
    if result.Error != nil {
        return result.Error
    }
    return nil

}

func (s MariaDBStore) Get(Id int) (*JobApplication, error) {
	var application JobApplication
	result := s.db.First(&application, Id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &application, nil
} 

func (s MariaDBStore) List() ([]JobApplication, error){
	var applications []JobApplication
	s.db.Find(&applications)
	return applications, nil
}


