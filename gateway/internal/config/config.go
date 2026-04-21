package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPConfig struct {
	Address string        `yaml:"address" env:"ADDRESS" env-default:":8080"`
	Timeout time.Duration `yaml:"timeout" env:"API_TIMEOUT" env-default:"5s"`
}

type Config struct {
	LogLevel       string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
	HTTPConfig     HTTPConfig `yaml:"api-server"`
	AuthAddress    string     `yaml:"auth_address" env:"AUTH_ADDRESS" env-default:"auth-service:8080"`
	UserAddress    string     `yaml:"user_address" env:"USER_ADDRESS" env-default:"user-service:8080"`
	ProjectAddress string     `yaml:"project_address" env:"PROJECT_ADDRESS" env-default:"project-service:8080"`
	ReportAddress  string     `yaml:"report_address" env:"REPORT_ADDRESS" env-default:"report-service:8080"`
}

func MustLoad(configPath string) *Config {
	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("Error loading config: %s", err)
	}

	return &config
}