package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel  string `yaml:"log_level" env:"LOG_LEVEL" env-default:"debug"`
	Address   string `yaml:"address" env:"PROJECT_ADDRESS" env-default:":8080"`
	DBAddress string `yaml:"db_address" env:"DB_ADDRESS" env-default:"postgres://postgres:password@project-postgres:5432/projectpb?sslmode=disable"`
}

func MustLoad(configPath string) Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return cfg
}
