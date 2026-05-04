package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel       string `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	Address        string `yaml:"address" env:"PAYMENT_CALENDAR_ADDRESS" env-default:":8080"`
	DBAddress      string `yaml:"db_address" env:"DB_ADDRESS" env-default:"postgres://postgres:password@payment-calendar-postgres:5432/paymentcalendardb?sslmode=disable"`
	ProjectAddress string `yaml:"project_address" env:"PROJECT_ADDRESS" env-default:"project-service:8080"`
}

func MustLoad(configPath string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config %q: %s", configPath, err)
	}
	return &cfg
}
