package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port         string
	RedisAddress string
	PostgresURI  string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &Config{
		Port:         viper.GetString("port"),
		RedisAddress: viper.GetString("redis_address"),
		PostgresURI:  viper.GetString("postgres_uri"),
	}

	return cfg, nil
}
