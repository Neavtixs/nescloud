package configs

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnv(filename ...string) {
	err := godotenv.Load(filename...)
	if err != nil {
		log.Println("No .env file found, using system env variables")
	}
}
