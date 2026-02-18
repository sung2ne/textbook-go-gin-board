package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

type Config struct {
	Level      string
	Pretty     bool
	TimeFormat string
}

func Init(cfg Config) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	if cfg.TimeFormat != "" {
		zerolog.TimeFieldFormat = cfg.TimeFormat
	} else {
		zerolog.TimeFieldFormat = time.RFC3339
	}

	var output io.Writer = os.Stdout
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	}

	Logger = zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Info() *zerolog.Event  { return Logger.Info() }
func Error() *zerolog.Event { return Logger.Error() }
func Debug() *zerolog.Event { return Logger.Debug() }
func Warn() *zerolog.Event  { return Logger.Warn() }
func Fatal() *zerolog.Event { return Logger.Fatal() }
