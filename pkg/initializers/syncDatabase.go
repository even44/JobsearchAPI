package initializers

import "github.com/even44/JobsearchAPI/pkg/models"

func SyncDatabase() {
	err := Db.AutoMigrate(
		&models.JobApplication{}, 
		&models.Company{}, 
		&models.Contact{},
		&models.User{},
	)
	if err != nil {
		logger.Fatal(err)
	}
}