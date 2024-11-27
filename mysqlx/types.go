package mysqlx

import (
	"time"

	gLogger "gorm.io/gorm/logger"
)

type Config struct {
	Addr         string
	Username     string
	Password     string
	DBName       string
	MaxOpenConns int
	MaxIdleConns int

	// logger, you can use log.Logger or logrus.Logger, etc...
	LogLevel gLogger.LogLevel

	// slow query threshold
	SlowThreshold time.Duration
}

func (c *Config) setDefault() {
	if c.LogLevel == 0 {
		c.LogLevel = gLogger.Info
	}
	if c.SlowThreshold == 0 {
		c.SlowThreshold = 100 * time.Millisecond
	}
}
