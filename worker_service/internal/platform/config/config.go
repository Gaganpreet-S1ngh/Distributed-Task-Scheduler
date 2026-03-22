package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GrpcServerAddress string
	GrpcServerPort    string
}

func LoadConfig() Config {
	if os.Getenv("ENV") != "prod" {
		err := godotenv.Load()

		if err != nil {
			log.Println("Warning: .env file not found or could not be loaded")
		}
	}

	cfg := Config{
		GrpcServerAddress: os.Getenv("COORDINATOR_SERVER_ADDRESS"),
		GrpcServerPort:    os.Getenv("GRPC_SERVER_PORT"),
	}

	return cfg

}
