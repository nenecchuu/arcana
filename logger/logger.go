package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type Config struct {
	ServiceName string
	Level       zerolog.Level
}

func InitGlobalLogger(cfg *Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(cfg.Level)

	// Write errors and fatal errors to Stderr
	errWriter := &ErrorLoggerWriter{zerolog.MultiLevelWriter(os.Stderr)}
	// Write debug, informational, and warnings to Stdout
	infoWriter := &InfoLoggerWriter{zerolog.MultiLevelWriter(os.Stdout)}

	lw := zerolog.MultiLevelWriter(errWriter, infoWriter)

	log.Logger = zerolog.New(lw).With().Timestamp().Caller().Str("service", cfg.ServiceName).Logger()
}

func SetGlobalLevel(lvl zerolog.Level) {
	zerolog.SetGlobalLevel(lvl)
}
