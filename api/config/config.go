package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_USERNAME         string
	DB_PASSWORD         string
	DB_HOST             string
	DB_PORT             string
	DB_NAME             string
	STORAGE_BUCKET_NAME string
	GOOGLE_PROJECT_ID   string
	KAFKA_TRACE_TOPIC   string
	KAFKA_BROKER        string
}

var envConfig = fetchEnv()

func fetchEnv() Config {
	_ = godotenv.Load()

	return Config{
		DB_USERNAME:         GetEnv("DB_USERNAME"),
		DB_PASSWORD:         GetEnv("DB_PASSWORD"),
		DB_HOST:             GetEnv("DB_HOST"),
		DB_PORT:             GetEnv("DB_PORT"),
		DB_NAME:             GetEnv("DB_NAME"),
		STORAGE_BUCKET_NAME: GetEnv("STORAGE_BUCKET_NAME"),
		GOOGLE_PROJECT_ID:   GetEnv("GOOGLE_PROJECT_ID"),
		KAFKA_TRACE_TOPIC:   GetEnv("KAFKA_TRACE_TOPIC"),
		KAFKA_BROKER:        GetEnv("KAFKA_BROKER"),
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	return value
}

func GetEnvConfig() Config {
	return envConfig
}
