package logger

import (
	"context"
	"github.com/oaxacos/vitacare/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	once             = sync.Once{}
	log              *zap.SugaredLogger
	loggerContextKey = "ctx_logger"
)

func setGlobalLogger(logger *zap.SugaredLogger) {
	log = logger
}

func New(config *config.Config) *zap.SugaredLogger {
	return NewLogger(config.Server.Debug, config.Server.Pretty)
}

func NewLogger(debugMode, prettyLog bool) *zap.SugaredLogger {
	//Define the encoder with color support
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // Human-readable time format
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, // Short file path
	}
	if prettyLog {
		// Adds color to logs
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Create a core with the level and colored output
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	logLevel := zapcore.InfoLevel
	if debugMode {
		logLevel = zapcore.DebugLevel
	}

	// Output to stdout
	consoleOutput := zapcore.Lock(os.Stdout)

	// Combine the encoder, level, and output into a core
	core := zapcore.NewCore(consoleEncoder, consoleOutput, logLevel)

	// Build the logger
	return zap.New(core, zap.AddCaller()).Sugar()

}

func GetGlobalLogger(params ...bool) *zap.SugaredLogger {
	once.Do(func() {
		var log *zap.SugaredLogger
		if len(params) == 0 {
			// default logger is in debug mode and pretty logs
			log = NewLogger(true, true)
		} else if len(params) == 1 {
			log = NewLogger(params[0], params[0])
		} else {
			log = NewLogger(params[0], params[1])
		}

		setGlobalLogger(log)
	})

	return log
}

func CloseLogger() {
	if log != nil {
		_ = log.Sync()
	}
}

func GetContextLogger(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		log := GetGlobalLogger()
		return log
	}
	if log, ok := ctx.Value(loggerContextKey).(*zap.SugaredLogger); ok {
		return log
	}
	return GetGlobalLogger()
}

func SetContextLogger(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}
