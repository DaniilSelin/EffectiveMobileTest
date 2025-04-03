package config

import (
	"fmt"
	"os"
	
	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Dbname        string `yaml:"dbname"`
	Sslmode       string `yaml:"sslmode"`
	Schema		  string `yaml:"schema"`
}

func (db DatabaseConfig) ConnString() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s&search_path=%s",
		db.User, db.Password, db.Host, db.Port, db.Dbname, db.Sslmode, db.Schema,
	)
}

type Config struct {
	DB DatabaseConfig `yaml:"database"`
	Server  ServerConfig `yaml:"server"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("could not decode config file: %v", err)
	}
	return &config, nil
}
