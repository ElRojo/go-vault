package utility

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch logLevel := os.Getenv("LOG_LEVEL"); logLevel {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("log level set to debug")
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
