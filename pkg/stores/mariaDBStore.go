package stores

import (
	"log"
	"os"

	"gorm.io/gorm"
)

type MariaDBStore struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewMariaDBStore(db *gorm.DB) *MariaDBStore {
	var logger = log.New(os.Stdout, "MariaDBStore: ", log.Ldate+log.Ltime+log.Lmsgprefix)
	logger.Println("Auto migrating database")

	return &MariaDBStore{
		db,
		logger,
	}
}
