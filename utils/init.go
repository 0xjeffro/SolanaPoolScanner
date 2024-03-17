package utils

import (
	"SolanaPoolScanner/db"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func checkEnv() {
	necessaryEnv := []string{
		"API_TOKEN", "DSN",
		"TICKER", "TOKEN_ADDR",
	}
	for _, env := range necessaryEnv {
		if os.Getenv(env) == "" {
			panic("Missing environment variable: " + env)
		}
	}
}

func initLogger() {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig = encoderConfig
	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func Init() {
	initLogger()
	checkEnv()
	db.CreateTable()
}
