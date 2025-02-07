package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var logger *log.Logger

func LoadEnvVariables() {
	logger = log.New(os.Stdout, "INIT: ", log.Ldate+log.Ltime+log.Lmsgprefix)
	logger.Println("Loading .env file")
	err := godotenv.Load()
	if err != nil {
		logger.Println("[ERROR] No .env file or error loading, skipping")
	}
}
