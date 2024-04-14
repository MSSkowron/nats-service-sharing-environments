package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServerAddress string `envconfig:"SERVER_ADDRESS" default:""`
	Username      string `envconfig:"USERNAME" default:""`
	Publish       bool   `envconfig:"PUBLISH" default:"true"`
	Subscribe     bool   `envconfig:"SUBSCRIBE" default:"true"`
	ListenSelf    bool   `envconfig:"LISTEN_SELF" default:"false"`
	Interactive   bool   `envconfig:"INTERACTIVE" default:"false"`
	Timeout       int    `envconfig:"TIMEOUT" default:"1000"`
	AddStream     bool   `envconfig:"ADD_STREAM" default:"false"`
	Subject       string `envconfig:"SUBJEST" default:""`
}

func Load() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("error loading configuration from environment variables: %w", err)
	}
	return cfg, nil
}
