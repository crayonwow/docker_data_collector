package logger

import (
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
		Level uint32 `yaml:"level"`
	}
)

func newLogger(c *Config) *logrus.Logger {
	logrus.Info(c)

	l := logrus.StandardLogger()
	l.SetLevel(logrus.Level(c.Level))
	l.SetReportCaller(true)
	return l
}
