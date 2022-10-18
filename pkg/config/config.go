package config

import (
	"errors"
	"fmt"
	"os"

	"docker_data_collector/pkg/logger"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		Logger logger.Config

		RawConfig []byte
	}
)

func NewConfig(params configIn) (*Config, error) {
	if params.ConfigPath == "" {
		return nil, errors.New("empty path")
	}
	b, err := os.ReadFile(params.ConfigPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	c := &Config{
		RawConfig: b,
	}
	return c, yaml.Unmarshal(b, c)
}

func LoggerConfig(config *Config) *logger.Config {
	return &config.Logger
}
