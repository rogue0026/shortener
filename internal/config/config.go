package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Shortener struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address      string        `yaml:"address"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

func MustLoad() Shortener {
	path := os.Getenv("CONFIG")
	if path == "" {
		panic("env CONFIG not defined")
	}
	_, err := os.Stat(path)
	if err != nil {
		panic("config file does not exist")
	}
	cfg := Shortener{}
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic(err)
	}

	return cfg
}
