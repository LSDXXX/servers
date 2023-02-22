package config

import lumberjack "gopkg.in/natefinch/lumberjack.v2"

// LogConfig log conf
type LogConfig struct {
	Output     lumberjack.Logger `yaml:"output"`
	Level      string            `yaml:"level" default:"debug"`
	WithCaller bool              `yaml:"with_caller" default:"true"`
	WithStdOut bool              `yaml:"with_std_out"`
	HiddenKey  bool
}
