package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	Params   string
}

var requiredEnvs = []string{
	"DATABASE_NAME",
	"DATABASE_USER",
	"DATABASE_PASSWORD",
	"DATABASE_HOST",
	"DATABASE_PORT",
}

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func requireEnvs() error {
	missingEnvs := make([]string, 0)
	for _, key := range requiredEnvs {
		if getEnv(key) == "" {
			missingEnvs = append(missingEnvs, key)
		}
	}

	if len(missingEnvs) > 0 {
		return fmt.Errorf("missing required env variables: %s", strings.Join(missingEnvs, ", "))
	}

	return nil
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, using environment instead")
	}

	if err := requireEnvs(); err != nil {
		return Config{}, err
	}

	return Config{
		DB: DBConfig{
			Name:     getEnv("DATABASE_NAME"),
			User:     getEnv("DATABASE_USER"),
			Password: getEnv("DATABASE_PASSWORD"),
			Host:     getEnv("DATABASE_HOST"),
			Port:     getEnv("DATABASE_PORT"),
			Params:   getEnv("DATABASE_PARAMS"),
		},
	}, nil
}
