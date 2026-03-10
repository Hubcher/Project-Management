package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Address string        `yaml:"address" env:"ADDRESS" env-default:"localhost:80"`
	Timeout time.Duration `yaml:"timeout" env:"API_TIMEOUT" env-default:"5s"`
}
type Config struct {
	LogLevel       string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	HTTPConfig     HTTPConfig `yaml:"api-server"`
	UserAddress    string     `yaml:"user_address" env:"USER_ADDRESS" env-default:"localhost:81"`
	ProjectAddress string     `yaml:"project_address" env:"PROJECT_ADDRESS" env-default:"localhost:82"`
}

func MustLoad(configPath string) Config {
	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	return config
}
