package logger

import "go.uber.org/zap/zapcore"

type Config struct {
	// Name Root Logger 的名称
	Name             string        `json:"name" yaml:"name" mapstructure:"name"`
	Level            zapcore.Level `json:"level" yaml:"level" mapstructure:"level"`
	StacktraceLevel  zapcore.Level `json:"stacktraceLevel" yaml:"stacktraceLevel" mapstructure:"stacktraceLevel"`
	Console          bool          `json:"console" yaml:"console" mapstructure:"console"`
	LumberjackConfig `json:",inline" yaml:",squash" mapstructure:",squash"`
	Fields           []zapcore.Field `json:"fields" yaml:"fields" mapstructure:"fields"`
}

type LumberjackConfig struct {
	FilePath string `json:"filePath" yaml:"filePath" mapstructure:"filePath"`
	MaxSize  int    `json:"maxSize" yaml:"maxSize" mapstructure:"maxSize"`
	MaxAge   int    `json:"maxAge" yaml:"maxAge" mapstructure:"maxAge"`
}

func Encoder(console bool) zapcore.Encoder {
	config := zapcore.EncoderConfig{
		NameKey:        "name",
		TimeKey:        "time",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeName:     zapcore.FullNameEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if console {
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
		return zapcore.NewConsoleEncoder(config)
	}

	return zapcore.NewJSONEncoder(config)
}

func (c *Config) setDefault() {
	if c.Name == "" {
		c.Name = "default"
	}
	if c.StacktraceLevel == 0 {
		c.StacktraceLevel = zapcore.DPanicLevel
	}
}
