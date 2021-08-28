package utils

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig loads config file from given path.
func LoadConfig(filepath string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filepath)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// GetConfigPath returns config path for local or docker environment.
func GetConfigPath(serviceName, configEnv string) string {
	if configEnv == "docker" {
		return fmt.Sprintf("./config/%s/config-docker", serviceName)
	}

	return fmt.Sprintf("./config/%s/config-local", serviceName)
}
