package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

// Структура конфига
type Config struct {
	Database struct {
		URL string `yaml:"url"`
	} `yaml:"database"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	JWT struct {
		Secret string `yaml:"secret"`
	} `yaml:"jwt"`
	SMTP struct {
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"smtp"`
}

var AppConfig *Config

// Загрузка и парсинг YAML-файла
func Load() *Config {
	data, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка чтения config.yaml: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Ошибка парсинга YAML: %v", err)
	}

	AppConfig = &cfg
	return AppConfig
}
