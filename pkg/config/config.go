package config

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"docker_data_collector/pkg/logger"

	"gopkg.in/yaml.v3"
)

type (
	Config struct {
		RawConfig []byte
		Logger    logger.Config
	}
)

func configPath() string {
	path := ""
	flag.StringVar(&path, "configPath", "config/config.yaml", "path to the config file")

	return path
}

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
