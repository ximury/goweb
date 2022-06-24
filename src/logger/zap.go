package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type zapLogger struct {
	sugaredLogger *zap.SugaredLogger
}

//默认编码
func defaultEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

//获取level
func getZapLevel(level string) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func newZapLogger(config Configuration, args ...interface{}) (Logger, error) {
	var (
		core  zapcore.Core
		cores []zapcore.Core
	)
	//输出到控制台，main函数初始化EnableConsole的值
	if config.EnableConsole {
		//获取日志文件级别
		level := getZapLevel(config.ConsoleLevel)
		//写入Stdout(标准输出)
		writer := zapcore.Lock(os.Stdout)
		//NewCore==定制Logger 1、Encode：编码器；2、日志写入哪里？使用AddSync()函数将打开的文件句柄传进去。
		core := zapcore.NewCore(defaultEncoder(config.ConsoleJSONFormat), writer, level)
		cores = append(cores, core)
	}
	//系统默认值--->输出到日志文件
	if config.EnableFile {
		//获取日志文件级别
		level := getZapLevel(config.FileLevel)
		//文件句柄
		writer := zapcore.AddSync(&lumberjack.Logger{
			//文件地址
			Filename: config.FileLocation,
			//最大空间 100MB
			MaxSize: config.MaxSize,
			//压缩
			Compress: config.Compress,
			//时长 28天
			MaxAge: config.MaxAge,
		})
		//NewCore==定制Logger 1、Encode：编码器；2、日志写入哪里？使用AddSync()函数将打开的文件句柄传进去。
		core := zapcore.NewCore(defaultEncoder(config.FileJSONFormat), writer, level)
		cores = append(cores, core)
	}

	// 用户初始化值--->输出到日志文件
	if config.UserDefine {
		level := getZapLevel(config.FileLevel)
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:  config.FileLocation,
			MaxSize:   config.MaxSize,
			Compress:  config.Compress,
			MaxAge:    config.MaxAge,
			LocalTime: config.LocalTime,
		})
		for _, arg := range args {
			switch arg.(type) {
			case zapcore.Encoder:
				core = zapcore.NewCore(arg.(zapcore.Encoder), writer, level)
			default:
				core = zapcore.NewCore(defaultEncoder(config.FileJSONFormat), writer, level)
			}
		}

		cores = append(cores, core)
	}
	combinedCore := zapcore.NewTee(cores...)

	// AddCallerSkip skips 2 number of callers, this is important else the file that gets
	// logged will always be the wrapped file. In our case zap.go
	logger := zap.New(combinedCore,
		zap.AddCallerSkip(2),
		zap.AddCaller(),
	).Sugar()

	return &zapLogger{
		sugaredLogger: logger,
	}, nil
}

func (l *zapLogger) Debugf(format string, args ...interface{}) {
	l.sugaredLogger.Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...interface{}) {
	l.sugaredLogger.Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...interface{}) {
	l.sugaredLogger.Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...interface{}) {
	l.sugaredLogger.Errorf(format, args...)
}

func (l *zapLogger) Fatalf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) Panicf(format string, args ...interface{}) {
	l.sugaredLogger.Fatalf(format, args...)
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	var f = make([]interface{}, 0)
	for k, v := range fields {
		f = append(f, k)
		f = append(f, v)
	}
	newLogger := l.sugaredLogger.With(f...)
	return &zapLogger{newLogger}
}
