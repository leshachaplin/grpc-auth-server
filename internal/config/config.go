package config

// if using go modules

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port             int    `env:"Port" envDefault:"1323"`
	GrpcPort         string `env:"GrpcPort" envDefault:"0.0.0.0:50051"`
	SecretKeyAuth    string `env:"SecretKeyAuth"`
	SecretKeyRefresh string `env:"SecretKeyRefresh"`
	PostgresConStr   string `env:"PostgresConStr"`
	MongoConStr      string `env:"MongoConStr"`
	SMTPConStr       string `env:"SMTPConStr"`
}

func NewConfig() *Config {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}
	return &cfg
}
