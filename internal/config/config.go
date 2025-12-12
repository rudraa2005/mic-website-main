package config

import (
	"log"
	"os"
)

type Config struct {
	Env           string
	Port          string
	DatabaseURL   string
	JWTSecret     string
	JWTExpMinutes string
	PGMaxConns    string
	RedisURL      string
	KafkaBrokers  string
	AdminEmail    string
}

func Load() *Config {
	cfg := &Config{
		Env:           os.Getenv("ENV"),
		Port:          os.Getenv("PORT"),
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		JWTExpMinutes: os.Getenv("JWT_EXP_MINUTES"),
		PGMaxConns:    os.Getenv("PG_MAX_CONNS"),
		RedisURL:      os.Getenv("REDIS_URL"),
		KafkaBrokers:  os.Getenv("KAFKA_BROKERS"),
		AdminEmail:    os.Getenv("ADMIN_EMAIL"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL missing")
	}

	return cfg
}


