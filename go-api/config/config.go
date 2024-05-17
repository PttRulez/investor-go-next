package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type PostgresConfig struct {
	DBName   string
	SSLMode  string
	Host     string
	Password string
	Port     string
	Username string
}

type Config struct {
	ApiHost         string
	ApiPort         int
	AllowedCors     []string
	JwtSecret       string
	Pg              PostgresConfig
	TokenAuthSecret string
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
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenAuthSecret := os.Getenv("TOKEN_AUTH_SECRET")

	return &Config{
		AllowedCors:     strings.Split(corsString, ","),
		ApiPort:         apiPort,
		ApiHost:         apiHost,
		JwtSecret:       jwtSecret,
		Pg:              pgConfig,
		TokenAuthSecret: tokenAuthSecret,
	}
}
