package logger

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(logFile string) *zap.Logger {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // Максимальный размер файла в MB
		MaxBackups: 5,  // Количество старых файлов
		MaxAge:     30, // Максимальный возраст файлов в днях
		Compress:   true,
	}

	writeSyncer := zapcore.AddSync(lumberjackLogger)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zapcore.DebugLevel,
	)

	return zap.New(core)
}