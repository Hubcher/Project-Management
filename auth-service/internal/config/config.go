package config

import (
    "log"
    "time"

    "github.com/ilyakaznacheev/cleanenv"
)

type JWT struct {
    Secret   string        `yaml:"secret"`
    Issuer   string        `yaml:"issuer"`
    Audience string        `yaml:"audience"`
    TTL      time.Duration `yaml:"ttl"`
}

type BootstrapAdmin struct {
    Enabled  bool   `yaml:"enabled" env:"BOOTSTRAP_ADMIN_ENABLED" env-default:"true"`
    Email    string `yaml:"email" env:"BOOTSTRAP_ADMIN_EMAIL" env-default:"admin@example.com"`
    Password string `yaml:"password" env:"BOOTSTRAP_ADMIN_PASSWORD" env-default:"Admin123!"`
}

type Config struct {
    LogLevel       string         `yaml:"log_level" env:"LOG_LEVEL" env-default:"DEBUG"`
    Address        string         `yaml:"address" env:"AUTH_ADDRESS" env-default:":50051"`
    Env            string         `yaml:"env" env:"ENV" env-default:"local"`
    DBAddress      string         `yaml:"db_address" env:"DB_ADDRESS" env-default:"postgres://postgres:password@auth-postgres:5432/authdb?sslmode=disable"`
    JWT            JWT            `yaml:"jwt"`
    BootstrapAdmin BootstrapAdmin `yaml:"bootstrap_admin"`
}

func MustLoad(configPath string) *Config {
    var cfg Config

    if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
        log.Fatalf("cannot read config %q: %s", configPath, err)
    }

    return &cfg
}
