package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel  string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address   string `yaml:"address" env:"USER_ADDRESS" env-default:":8080"`
	DBAddress string `yaml:"db_address" env:"DB_ADDRESS" env-default:"postgres://postgres:password@user-postgres:5432/userdb?sslmode=disable"`
}

func MustLoad(configPath string) Config {
	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return config
}
