package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string           `yaml:"env" env-default:"local"`
	StoragePath string           `yaml:"storage_path" env-required:"true"`
	HttpServer  HttpServerConfig `yaml:"http_server" env-required:"true"`
}

type HttpServerConfig struct {
	Adress        string        `yaml:"adress" env-default:"localhost:8080"`
	Timeout       time.Duration `yaml:"timeout" env-default:"4s"`
	Iddle_timeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

func MustLoadByPath(configPath string) *Config {
	// Проверка наличия файла
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}
