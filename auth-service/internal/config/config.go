package config

import (
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type GRPCConfig struct {
	Address string        `yaml:"address" env:"GRPC_ADDRESS" env-default:":50052"`
	Timeout time.Duration `yaml:"timeout" env:"GRPC_TIMEOUT" env-default:"5s"`
}

type JWTConfig struct {
	Secret string        `yaml:"secret" env:"JWT_SECRET" env-required:"true"`
	TTL    time.Duration `yaml:"ttl" env:"JWT_TTL" env-default:"15m"`
}

type Config struct {
	LogLevel           string     `yaml:"log_level" env:"LOG_LEVEL" env-default:"INFO"`
	GRPC               GRPCConfig `yaml:"grpc"`
	JWT                JWTConfig  `yaml:"jwt"`
	UserServiceAddress string     `yaml:"user_service_address" env:"USER_SERVICE_ADDRESS" env-default:"user-service:50051"`
}

func MustLoad(configPath string) *Config {
	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	return &cfg
}

//type Config struct {
//	LogLevel  string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"debug"`
//	Address   string        `yaml:"address" env:"PROJECT_ADDRESS" env-default:":8080"`
//	Env       string        `yaml:"env" env:"local" env-default:"local"`
//	DBAddress string        `yaml:"db_address" env:"DB_ADDRESS" env-default:"postgres://postgres:password@auth-postgres:5432/projectpb?sslmode=disable"`
//	TokenTTL  time.Duration `yaml:"token_ttl" env:"TOKEN_TTL" env-required:"true"`
//}
//
//func MustLoad(configPath string) *Config {
//	var cfg Config
//
//	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
//		log.Fatalf("cannot read config %q: %s", configPath, err)
//	}
//	return &cfg
//}
