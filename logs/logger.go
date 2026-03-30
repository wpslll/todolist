package logs

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Newlogger(logLevel string) (*zap.Logger, func() error, error) {
	lvl := zap.NewAtomicLevel()
	if err := lvl.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, nil, err
	}
	if err := os.MkdirAll("logFiles", 0755); err != nil {
		return nil, nil, err
	}
	timestamp := time.Now().UTC()
	logFilePath := filepath.Join("logFiles", fmt.Sprintf("%s.log", timestamp))
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, nil, err
	}
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(logFile), lvl),
	)
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	return logger, logFile.Close, nil
}
