package eigenlayer

import (
	"log/slog"
	"os"

	eigensdkLogger "github.com/Layr-Labs/eigensdk-go/logging"
)

func GetLogger(verbose bool) eigensdkLogger.Logger {
	loggerOptions := &eigensdkLogger.SLoggerOptions{
		Level: slog.LevelInfo,
	}
	if verbose {
		loggerOptions = &eigensdkLogger.SLoggerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}
	}
	logger := eigensdkLogger.NewTextSLogger(os.Stdout, loggerOptions)
	return logger
}
