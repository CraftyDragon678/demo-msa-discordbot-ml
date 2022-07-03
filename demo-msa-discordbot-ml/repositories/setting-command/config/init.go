package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Token string
)

func init() {
	path, _ := os.Getwd()
	if err := godotenv.Load(path + "/.env"); err != nil {
		log.Fatalf("Error: %v", err)
	}

	Token = os.Getenv("BOT_TOKEN")
}
