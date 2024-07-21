package logger

import (
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/ciscapello/notification_service/application/config"
	"go.uber.org/zap/zapcore"
)

var (
	logger *slog.Logger
	once   sync.Once
)

func GetLogger(config *config.Config) *slog.Logger {
	once.Do(func() {
		logger = InitLogger(config.LogPath)
	})
	return logger
}

func createLogDirectoryIfNotExists(logPath string) {
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		err = os.Mkdir(logPath, 0777)
		if err != nil {
			panic(err)
		}
	}
}

func getLogWriter(logPath string) zapcore.WriteSyncer {
	file, err := os.OpenFile(logPath+"/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	multiwriter := io.MultiWriter(file, os.Stdout)

	if err != nil {
		panic(err)
	}
	return zapcore.AddSync(multiwriter)
}

func InitLogger(logPath string) *slog.Logger {
	createLogDirectoryIfNotExists(logPath)
	writer := getLogWriter(logPath)
	handler := slog.NewJSONHandler(writer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	return logger
}
