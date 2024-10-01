package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pttrulez/investor-go/internal/api"
)

type APIConfig struct {
	APIHost         string
	APIPort         int
	AllowedCors     []string
	TokenAuthSecret string
}

type PostgresConfig struct {
	DBName   string
	SSLMode  string
	Host     string
	Password string
	Port     string
	Username string
}

type Config struct {
	API api.Config
	Pg  PostgresConfig
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	pgConfig := PostgresConfig{
		DBName:   os.Getenv("PG_DB_NAME"),
		Host:     os.Getenv("PG_HOST"),
		Password: os.Getenv("PG_PASSWORD"),
		Port:     os.Getenv("PG_PORT"),
		SSLMode:  os.Getenv("PG_SSLMODE"),
		Username: os.Getenv("PG_USERNAME"),
	}

	apiPort, _ := strconv.Atoi(os.Getenv("GO_API_PORT"))
	apiHost := os.Getenv("GO_API_HOST")
	corsString := os.Getenv("CORS_ALLOWED_ORIGINS")
	allowedCors := strings.Split(corsString, ",")
	tokenAuthSecret := os.Getenv("TOKEN_AUTH_SECRET")

	apiConfig := api.Config{
		AllowedCors:     allowedCors,
		APIPort:         apiPort,
		APIHost:         apiHost,
		TokenAuthSecret: tokenAuthSecret,
	}

	return &Config{
		API: apiConfig,
		Pg:  pgConfig,
	}
}
