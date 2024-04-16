package helpers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(env string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	data := os.Getenv(env)

	return data

}
