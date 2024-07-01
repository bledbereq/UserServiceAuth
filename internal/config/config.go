package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string           `yaml:"env" env-default:"local"`
	StoragePath string           `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration    `yaml:"token_ttl" env-required:"true"`
	GRPS        GRPSconfig       `yaml:"grpc" env-required:"true"`
	HttpServer  HttpServerConfig `yaml:"http_server" env-required:"true"`
}
type GRPSconfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type HttpServerConfig struct {
	Adress        string        `yaml:"adress" env-default:"localhost:8080"`
	Timeout       time.Duration `yaml:"timeout" env-default:"4s"`
	Iddle_timeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("config path is empty: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
