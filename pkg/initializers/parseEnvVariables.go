package initializers

import (
	"os"
	"strconv"
)

// Api Envs
var ApiPort int = 3001
var ApiTrustedOrigin string = ""
var ApiSecret string = "setthisasasecretenv"

// Db Envs
var DbURL = ""
var DbUser = ""
var DbPassword = ""
var DbPort = 0

func ParseEnvVariables() {
	logger.Println("Getting API env variables")
	var temp string

	// Should look like "6001" not "sixthousandandone"
	temp = os.Getenv("API_PORT")
	if temp != "" {
		var err error
		ApiPort, err = strconv.Atoi(temp)
		if err != nil {
			logger.Println("[ERROR] Could not convert API_PORT to int")
			panic(err)
		}
	}

	// Should look like "http://ip:port" or "https://domain.example"
	temp = os.Getenv("TRUSTED_ORIGIN")
	if temp != "" {
		ApiTrustedOrigin = temp
	}

	// Should look like "192.168.0.20" not "sixthousandandone"
	temp = os.Getenv("DB_URL")
	if temp != "" {
		DbURL = temp
	}

	temp = os.Getenv("DB_USER")
	if temp != "" {
		DbUser = temp
	}

	temp = os.Getenv("DB_PASSWORD")
	if temp != "" {
		DbPassword = temp
	}

	temp = os.Getenv("DB_PORT")
	if temp != "" {
		var err error
		DbPort, err = strconv.Atoi(temp)
		if err != nil {
			logger.Fatal("[ERROR] Could not convert DB_PORT to int")
			panic(err)
		}
	}

	temp = os.Getenv("DB_PORT")
	if temp != "" {
		ApiSecret = temp
	}
}
