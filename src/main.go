package main

import (
	"logger"
	"main/config"
	"main/service"
)

func main() {
	service.StartApi()
}

func init() {
	configBase, err := config.GetChannelConfig()
	if err != nil {
		logger.Fatalf("Get channel config failed! err: %v", err)
	}
	//为日志指定参数
	configInit := logger.Configuration{
		EnableConsole:     configBase.Log.EnableConsole,
		ConsoleJSONFormat: configBase.Log.ConsoleJSONFormat,
		ConsoleLevel:      logger.GetLevel(configBase.Log.ConsoleLevel),
		EnableFile:        configBase.Log.EnableFile,
		FileJSONFormat:    configBase.Log.FileJSONFormat,
		FileLevel:         logger.GetLevel(configBase.Log.FileLevel),
		FileLocation:      configBase.Log.FileLocation,
		MaxAge:            configBase.Log.MaxAge,
		MaxSize:           configBase.Log.MaxSize,
		Compress:          configBase.Log.Compress,
	}
	err = logger.InitGlobalLogger(configInit, logger.InstanceZapLogger)
	if err != nil {
		logger.Fatalf("Could not instantiate log! err: %v", err)
	}
}
