package logger

import (
	"context"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/kirychukyurii/wasker-directory/internal/config"
)

type LoggerCtxKey struct{}

type Logger struct {
	zerolog.Logger
}

func New(config config.Config) Logger {
	var log zerolog.Logger
	var output io.Writer = zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	log = zerolog.New(output).
		Level(parseLogLevel(config.Log.Level)).
		With().
		Timestamp().
		Caller().
		Logger()

	return Logger{log}
}

func parseLogLevel(level string) zerolog.Level {
	switch strings.ToLower(level) {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *Logger) FromContext(ctx context.Context) *Logger {
	logger := ctx.Value(LoggerCtxKey{})

	if logger != nil {
		l.Logger = logger.(zerolog.Logger)
	}

	return l
}
