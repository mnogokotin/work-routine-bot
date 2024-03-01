package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type (
	Config struct {
		Env   string `yaml:"env"`
		Tg    Tg     `yaml:"tg"`
		Bot   Bot    `yaml:"bot"`
		Mongo Mongo  `yaml:"mongo"`
	}

	Tg struct {
		Host  string `yaml:"host" env:"TG_HOST"`
		Token string `yaml:"token" env:"TG_TOKEN"`
	}

	Bot struct {
		BatchSize int `yaml:"batch_size" env:"BOT_BATCH_SIZE"`
	}

	Mongo struct {
		Uri            string        `yaml:"uri" env:"MONGO_URI"`
		ConnectTimeout time.Duration `yaml:"connect_timeout" env:"MONGO_CONNECT_TIMEOUT"`
		DbName         string        `yaml:"db_name" env:"MONGO_DB_NAME"`
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
