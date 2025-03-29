package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func createLogger(env string) *zap.Logger{

	stdout:=zapcore.AddSync(os.Stdout)

	log_file := zapcore.AddSync(&lumberjack.Logger{
		Filename: "logs/app.log",
		MaxSize: 2,
		MaxBackups: 3,
		MaxAge: 7,
	})

	// default setup => production 

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	
	switch env{
	case "DEVELOPMENT": 
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "PRODUCTION": 
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

    Cfg := zap.NewProductionEncoderConfig()
    Cfg.TimeKey = "timestamp"
    Cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	

	consoleCfg := Cfg
	consoleCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	fileCfg := Cfg

	consoleEncoder := zapcore.NewConsoleEncoder(consoleCfg)
    fileEncoder := zapcore.NewJSONEncoder(fileCfg)

    core := zapcore.NewTee(
        zapcore.NewCore(consoleEncoder, stdout, level),
        zapcore.NewCore(fileEncoder, log_file, level),
    )

    return zap.New(core)

}