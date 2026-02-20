package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL                     string
	JWTSecret                       string
	TokenExpireHours                int
	APIKeyIntelligenceIO            string
	Recommendation_Service_BASE_URL string
}

func NewConfig() Config {
	if err := godotenv.Load("../../.env"); err != nil {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}

	return Config{
		DatabaseURL:                     os.Getenv("DATABASE_URL"),
		JWTSecret:                       os.Getenv("JWT_SECRET"),
		TokenExpireHours:                24,
		APIKeyIntelligenceIO:            os.Getenv("API_KEY_INTELLIGENCE_IO"),
		Recommendation_Service_BASE_URL: "http://localhost:8000/api",
	}
}

var Cfg = NewConfig()
