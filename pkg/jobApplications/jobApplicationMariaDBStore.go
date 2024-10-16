package jobApplications

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MariaDBStore struct {
	db *gorm.DB
}

func NewMariaDBStore() *MariaDBStore{
	db, err := connectToMariaDB()
	if err != nil {
		log.Fatal(err)
	}
	
	err = db.AutoMigrate(&JobApplication{})
	if err != nil {
		log.Fatal(err)
	}

	return &MariaDBStore{
		db,
	}
}

func connectToMariaDB() (*gorm.DB, error) {
    dsn := "root:superroot@tcp(db:3306)/jobsearchdb?charset=utf8mb4&parseTime=True&loc=Local"
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

func (s MariaDBStore) Update(id int, application JobApplication) error {
	result := s.db.Save(application)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (s MariaDBStore) Remove(id int) error {
	
	
	application, err := s.Get(id)
	if err != nil {
		return err
	}
	result := s.db.Delete(application)
	if result.Error != nil {
		return result.Error
	}
	return nil
	
}


