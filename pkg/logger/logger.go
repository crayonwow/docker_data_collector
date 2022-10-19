package logger

import (
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Level uint32 `yaml:"level"`
	}
)

func NewLogger(c *Config) {
	l := logrus.StandardLogger()
	l.SetLevel(logrus.Level(c.Level))
	l.SetReportCaller(true)
}
