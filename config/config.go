package config

import (
	"os"

	"github.com/w1png/go-htmx-ecommerce-template/errors"
)

var ConfigInstance *Config

type Config struct {
	Port      string
	JWTSecret string

	StorageType string

	SqlitePath string
}

func InitConfig() error {
	var ok bool
	config := &Config{}

	if config.Port, ok = os.LookupEnv("PORT"); !ok {
		return errors.NewEnvironmentVariableNotFoundError("PORT")
	}
	if config.JWTSecret, ok = os.LookupEnv("JWT_SECRET"); !ok {
		return errors.NewEnvironmentVariableNotFoundError("JWT_SECRET")
	}

	if config.StorageType, ok = os.LookupEnv("STORAGE_TYPE"); !ok {
		return errors.NewEnvironmentVariableNotFoundError("STORAGE_TYPE")
	}

	if config.SqlitePath, ok = os.LookupEnv("SQLITE_PATH"); !ok {
		return errors.NewEnvironmentVariableNotFoundError("SQLITE_PATH")
	}

	ConfigInstance = config

	return nil
}
