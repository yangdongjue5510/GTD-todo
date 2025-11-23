package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	DBMaxOpenConnection     int
	DBMaxIdleConnection     int
	DBConnectionMaxLifeTime int
}

func Load() *Config {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("DB Config load failed. %v", err)
	}
	return &Config{
		DBHost:                  os.Getenv("DB_HOST"),
		DBPort:                  os.Getenv("DB_PORT"),
		DBUser:                  os.Getenv("DB_USER"),
		DBPassword:              os.Getenv("DB_PASSWORD"),
		DBName:                  os.Getenv("DB_NAME"),
		DBMaxOpenConnection:     mustAtoi(os.Getenv("DB_MAX_OPEN_CONNECTION")),
		DBMaxIdleConnection:     mustAtoi(os.Getenv("DB_MAX_IDLE_CONNECTION")),
		DBConnectionMaxLifeTime: mustAtoi(os.Getenv("DB_CONNECTION_MAX_LIFE_TIME")),
	}
}

func mustAtoi(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("%s must be a valid integer: %v", value, err)
	}
	return i
}
