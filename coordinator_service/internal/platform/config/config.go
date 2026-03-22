package config

import (
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpServerPort string
	GrpcServerPort string
	DatabaseDSN    string
}

func LoadConfig() Config {
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()

		if err != nil {
			log.Println("Warning: .env file not found or could not be loaded")
		}
	}

	u := &url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")),
		Host:     os.Getenv("POSTGRES_HOST") + ":" + os.Getenv("POSTGRES_PORT"),
		Path:     os.Getenv("POSTGRES_DB"),
		RawQuery: "sslmode=disable",
	}

	dsn := u.String()

	cfg := Config{
		HttpServerPort: os.Getenv("HTTP_SERVER_PORT"),
		GrpcServerPort: os.Getenv("GRPC_SERVER_PORT"),
		DatabaseDSN:    dsn,
	}

	return cfg

}
