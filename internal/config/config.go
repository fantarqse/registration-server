package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	TokenSecretKey string `mapstructure:"TOKEN_SECRET_KEY"`

	DatabaseDriverName     string `mapstructure:"DATABASE_DRIVER_NAME"`
	DatabaseDataSourceName string `mapstructure:"DATABASE_DATA_SOURCE_NAME"`
}

func NewConfig(path, name string) (config *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return config, err
	}
	err = viper.Unmarshal(&config)
	log.Printf("New config: %v\n", config)
	return config, err
}
