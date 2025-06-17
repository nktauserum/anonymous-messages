package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

type Config struct {
	Telegram struct {
		Token string `yaml:"token,omitempty"`
		Admin int64  `yaml:"admin"`
	} `yaml:"telegram"`
	Port int64 `yaml:"port"`
	TLS  struct {
		Certificate string `yaml:"crt"`
		Key         string `yaml:"key"`
	} `yaml:"tls"`
}

var (
	once   sync.Once
	config *Config
)

func MustLoadConfig() *Config {
	once.Do(func() {
		data, err := os.ReadFile("config/config.yml")
		if err != nil {
			log.Panicf("error opening config: %s", err.Error())
		}

		config = new(Config)
		err = yaml.Unmarshal(data, &config)
		if err != nil {
			log.Panicf("error parsing config: %s", err.Error())
		}
	})

	return config
}
