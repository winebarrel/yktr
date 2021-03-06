package yktr

import (
	"fmt"
	"time"
)

type Config struct {
	Addr     string `default:"127.0.0.1"`
	Port     int    `default:"8080"`
	Team     string
	Token    string
	CacheTTL time.Duration `toml:"cache_ttl" default:"0"`
	PerPage  int           `toml:"per_page" default:"5"`
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
