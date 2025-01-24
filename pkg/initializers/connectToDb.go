package initializers

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectToMariaDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/jobsearchdb?charset=utf8mb4&parseTime=True&loc=Local",
		DbUser, DbPassword, DbURL, DbPort)
	logger.Printf("Connection string: %s", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	Db = db
}
