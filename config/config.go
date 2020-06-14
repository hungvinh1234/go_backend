package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var config *Config

// Config .
type Config struct {
	MySigningKey string `envconfig:"MYSIGNINGKEY"`
	Port         string `envconfig:"PORT"`
	Debug        bool   `envconfig:"DEBUG"`

	Postgres struct {
		Host   string `envconfig:"POSTGRES_HOST"`
		Port   string `envconfig:"POSTGRES_PORT"`
		DBName string `envconfig:"POSTGRES_DBNAME"`
		User   string `envconfig:"POSTGRES_USER"`
		Pass   string `envconfig:"POSTGRES_PASS"`
		Schema string `envconfig:"POSTGRES_SCHEMA"`
	}
}

func init() {
	config = &Config{}

	// read from env
	godotenv.Load()
	err := envconfig.Process("", config)
	if err != nil {
		panic(fmt.Sprintf("Failed to decode config env: %v", err))
	}

	// default value
	if len(config.Port) == 0 {
		config.Port = "3000"
	}
}

// GetConfig .
func GetConfig() *Config {
	return config
}
