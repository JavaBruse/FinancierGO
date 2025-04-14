package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	URL string `yaml:"url"`
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

type JWTConfig struct {
	Secret string `yaml:"secret"`
}

type SMTPConfig struct {
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
	SMTP     SMTPConfig     `yaml:"smtp"`
}

func Load() (*Config, error) {
	// Читаем конфигурационный файл
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	// Создаем конфигурацию
	cfg := &Config{}

	// Парсим YAML
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
