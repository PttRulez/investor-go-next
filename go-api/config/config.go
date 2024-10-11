package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	api "github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/http-server"
)

type Postgres struct {
	DBName   string
	SSLMode  string
	Host     string
	Password string
	Port     string
	Username string
}

type Redis struct {
	Address  string
	Password string
	DB       int
}

type Config struct {
	API          api.Config
	Pg           Postgres
	TgClientPort string
	Redis        Redis
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	pgConfig := Postgres{
		DBName:   os.Getenv("PG_DB_NAME"),
		Host:     os.Getenv("PG_HOST"),
		Password: os.Getenv("PG_PASSWORD"),
		Port:     os.Getenv("PG_PORT"),
		SSLMode:  os.Getenv("PG_SSLMODE"),
		Username: os.Getenv("PG_USERNAME"),
	}

	apiPort, err := strconv.Atoi(os.Getenv("GO_API_PORT"))
	if err != nil {
		panic(err)
	}

	apiConfig := api.Config{
		AllowedCors:     strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ","),
		APIPort:         apiPort,
		APIHost:         os.Getenv("GO_API_HOST"),
		TokenAuthSecret: os.Getenv("TOKEN_AUTH_SECRET"),
	}

	redisConfig := Redis{
		Address:  os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}

	return &Config{
		API:          apiConfig,
		Pg:           pgConfig,
		Redis:        redisConfig,
		TgClientPort: os.Getenv("TG_CLIENT_PORT"),
	}
}
