package logger

import (
	"errors"
	"go.uber.org/zap/zapcore"
)

// A global variable so that log functions can be directly accessed
var log Logger

// The Level is a logging priority. Higher levels are more important.
type Level int8

// Fields Type to pass when we want to call WithFields for structured logging
type Fields map[string]interface{}

const (
	//Debug is for verbose message
	Debug = "debug"
	//Info is default log level
	Info = "info"
	//Warn is for logging messages about possible issues
	Warn = "warn"
	//Error is for logging errors
	Error = "error"
	//Fatal is for logging fatal messages. The system will stop after the message is logged.
	Fatal = "fatal"
)

const (
	InstanceZapLogger int = iota
	InstanceLogrusLogger
)

var (
	errInvalidLoggerInstance = errors.New("invalid logger instance")
)

//Logger it is our contract for the logger
type Logger interface {
	Debugf(format string, args ...interface{})

	Infof(format string, args ...interface{})

	Warnf(format string, args ...interface{})

	Errorf(format string, args ...interface{})

	Fatalf(format string, args ...interface{})

	Panicf(format string, args ...interface{})
}

// Configuration stores the config for the logger
// For some loggers there can only be one level across writers, for such the level of Console is picked by default
type Configuration struct {
	EnableConsole     bool
	ConsoleJSONFormat bool
	ConsoleLevel      string
	EnableFile        bool
	FileJSONFormat    bool
	FileLevel         string
	FileLocation      string
	MaxAge            int
	MaxSize           int
	Compress          bool
	LocalTime         bool
	UserDefine        bool
}

//newLogger returns an instance of logger
func newLogger(config Configuration, loggerInstance int, args ...interface{}) (Logger, error) {
	var (
		err    error
		logger Logger
	)

	if args == nil {
		args = append(args, nil)
	}

	switch loggerInstance {
	case InstanceZapLogger:
		for _, arg := range args {
			switch arg.(type) {
			case zapcore.Encoder:
				logger, err = newZapLogger(config, arg.(zapcore.Encoder))
			default:
				logger, err = newZapLogger(config)
			}
		}

		if err != nil {
			return nil, err
		}
		return logger, nil

	case InstanceLogrusLogger:
		return nil, nil

	default:
		return nil, errInvalidLoggerInstance
	}
}

// InitGlobalLogger return a globally unique logger
func InitGlobalLogger(config Configuration, loggerInstance int, args ...interface{}) error {
	logger, err := newLogger(config, loggerInstance, args...)
	if err != nil {
		return err
	}
	log = logger
	return nil
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	log.Panicf(format, args...)
}

func GetLevel(level string) string {
	switch level {
	case "Debug":
		return Debug
	case "Info":
		return Info
	case "Warn":
		return Warn
	case "Error":
		return Error
	case "Fatal":
		return Fatal
	default:
		return Debug
	}
}
