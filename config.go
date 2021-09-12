package yktr

import (
	"fmt"
)

type Config struct {
	Addr  string `default:"127.0.0.1"`
	Port  int    `default:"8080"`
	Team  string
	Token string
}

func (cfg *Config) Validate() error {
	if cfg.Team == "" {
		return fmt.Errorf("config error: 'team' is required")
	}

	if cfg.Token == "" {
		return fmt.Errorf("config error: 'token' is required")
	}

	return nil
}
