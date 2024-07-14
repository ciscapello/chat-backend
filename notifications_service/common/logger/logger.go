package logger

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/ciscapello/notification_service/application/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

func GetLogger(config *config.Config) *zap.Logger {
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

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
}

func InitLogger(logPath string) *zap.Logger {
	createLogDirectoryIfNotExists(logPath)
	writer := getLogWriter(logPath)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writer, zapcore.DebugLevel)
	logger := zap.New(core)

	return logger
}
