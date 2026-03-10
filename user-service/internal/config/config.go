package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel  string `yaml:"logLevel" env:"LOG_LEVEL" env-default:"debug"`
	Address   string `yaml:"address" env:"USER_ADDRESS" env-default:"localhost:80"`
	DBAddress string `yaml:"db_address" env:"DB_ADDRESS" env-default:"localhost:82"`
}

func MustLoad(configPath string) Config {
	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return config
}
