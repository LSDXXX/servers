package config

import "github.com/LSDXXX/libs/config"

var (
	conf *Config
)

// Config description
type Config struct {
	Common config.Config    `yaml:"common"`
	GinLog config.LogConfig `yaml:"gin_log"`
	Logic  LogicConfig      `yaml:"logic"`
}

// LogicConfig description
type LogicConfig struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Proxy    string `yaml:"proxy"`
}

// ServerConfig description
// @return *Config
func ServerConfig() *Config {
	return conf
}

// SetServerConfig description
// @param c
func SetServerConfig(c *Config) {
	conf = c
}
