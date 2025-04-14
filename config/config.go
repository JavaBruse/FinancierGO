package config

import (
	"os"
)

type Config struct {
	DBUrl      string
	JWTSecret  string
	ServerPort string
	SMTPUser   string
	SMTPPass   string
}

func Load() *Config {
	return &Config{
		DBUrl:      os.Getenv("DB_URL"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("8189"),
		SMTPUser:   os.Getenv("SMTP_USER"),
		SMTPPass:   os.Getenv("SMTP_PASS"),
	}
}
