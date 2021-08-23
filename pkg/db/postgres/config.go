package postgres

import (
	"log"

	"github.com/grigagod/compresso/pkg/utils"
)

// Config stores postgres connection config.
type Config struct {
	Host     string
	Port     string
	User     string
	DbName   string
	Password string
	Driver   string
}

// LoadConfig loads postgres config from given filepath.
func LoadConfig(filepath string) (*Config, error) {
	v, err := utils.LoadConfig(filepath)
	if err != nil {
		log.Printf("unable to load  postgres config, %v", err)
	}

	var c Config

	err = v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil

}
