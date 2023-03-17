package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitEnvironment() {
	// loads environment variables
	err := godotenv.Load()
	if err != nil && os.Getenv("GIN_MODE") == "debug" {
		log.Fatal("Error loading .env file")
	}

}
