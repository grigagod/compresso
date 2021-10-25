package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig loads config into dst from given path.
func LoadConfig(dst interface{}, filepath string) error {
	v, err := PreloadConfig(filepath)
	if err != nil {
		return err
	}

	err = v.Unmarshal(dst)
	return err
}

func PreloadConfig(filepath string) (*viper.Viper, error) {
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
	if configEnv == "" {
		return fmt.Sprintf("./configs/%s/config-local", serviceName)
	}

	return fmt.Sprintf("./configs/%s/config-%s", serviceName, configEnv)
}
