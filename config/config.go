package config

import (
	"github.com/jinzhu/configor"
)

// Config - Application configuration
type Config struct {
	Log   string `yaml:"log" default:"" env:"LOG_PATH"`
	Debug bool   `yaml:"debug" default:"false" env:"DEBUG"`

	HTTP struct {
		Binding          string   `yaml:"binding" default:"localhost:8080" env:"HTTP_BINDING"`
		EndpointPath     string   `yaml:"endpoint_path" default:"/mcp" env:"HTTP_ENDPOINT_PATH"`
		HeartbeatSeconds int      `yaml:"heartbeat_seconds" default:"30" env:"HTTP_HEARTBEAT_SECONDS"`
		AuthToken        string   `yaml:"auth_token" default:"" env:"HTTP_AUTH_TOKEN"`
		AllowedOrigins   []string `yaml:"allowed_origins" env:"HTTP_ALLOWED_ORIGINS"`
	} `yaml:"http"`

	Greeting struct {
		DefaultMessage string `yaml:"default_message" default:"Hello!" env:"GREETING_DEFAULT_MESSAGE"`
	} `yaml:"greeting"`
}

// LoadConfig - Load configuration file
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	err := configor.New(&configor.Config{
		Debug:      false,
		Verbose:    false,
		Silent:     true,
		AutoReload: false,
	}).Load(cfg, path)
	return cfg, err
}
