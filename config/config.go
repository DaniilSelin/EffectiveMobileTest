package config

import (
	"fmt"
	"log"
	
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Host string `env:"SERVER_HOST"`
	Port int    `env:"SERVER_PORT"`
}

type DatabaseConfig struct {
    Host     string `env:"DB_HOST"`
    Port     int    `env:"DB_PORT"`
    User     string `env:"DB_USER"`
    Password string `env:"DB_PASSWORD"`
    Dbname   string `env:"DB_NAME"`
    Sslmode  string `env:"DB_SSLMODE"`
    Schema   string `env:"DB_SCHEMA"`
}

func (db DatabaseConfig) ConnString() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		db.User, db.Password, db.Host, db.Port, db.Dbname, db.Sslmode, db.Schema,
	)
}

type Config struct {
	DB DatabaseConfig
	Server  ServerConfig 
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
	}

	var config Config
	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("could not parse environment variables: %v", err)
	}

	log.Printf("Loaded config: %+v", config)
	return &config, nil
}