package dynamicrouter

import (
	"strings"
)

type Config struct {
	connection string
	UserName   string
	Password   string
}

func NewConfig(connection string, userName string, password string) *Config {
	return &Config{connection:connection, UserName:userName, Password:password}
}

func (config *Config) Address() string {
	index := strings.IndexAny(config.connection, "/")
	return string([]byte(config.connection)[0:index])
}

func (config *Config) Path() string {
	index := strings.IndexAny(config.connection, "/")
	return string([]byte(config.connection)[index:])
}

