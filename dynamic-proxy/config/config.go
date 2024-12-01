package config

// import (
// 	"github.com/spf13/viper"
// )

type PostgresAddr struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

type Config struct {
	Host          string
	Port          string
	TargetBackend string
	RedisAddr     string
	PostgresAddr  PostgresAddr
	AuthToken     string
}

// func LoadConfig() (*Config, error) {
// 	viper.SetConfigName("config")
// 	viper.SetConfigType("json")
// 	viper.AddConfigPath(".")

// 	if err := viper.ReadInConfig(); err != nil {
// 		return nil, err
// 	}

// 	cfg := &Config{
// 		Port:         viper.GetString("port"),
// 		RedisAddress: viper.GetString("redis_address"),
// 		PostgresURI:  viper.GetString("postgres_uri"),
// 	}

// 	return cfg, nil
// }
