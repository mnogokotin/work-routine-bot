package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
)

type (
	Config struct {
		Env      string   `yaml:"env"`
		Tg       Tg       `yaml:"tg"`
		Postgres Postgres `yaml:"postgres"`
		Rabbitmq Rabbitmq `yaml:"rabbitmq"`
	}

	Tg struct {
		Token string `yaml:"token" env:"TG_TOKEN"`
	}

	Postgres struct {
		Uri string `yaml:"url" env:"PG_URI"`
	}

	Rabbitmq struct {
		Uri string `yaml:"url" env:"RABBITMQ_URI"`
	}
)

func New() *Config {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/main.yml", cfg)
	if err != nil {
		panic("can't read config: " + err.Error())
	}

	err = cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		panic("can't read config: " + err.Error())
	}

	return cfg
}
