package utility

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch logLevel := strings.ToLower(os.Getenv("LOG_LEVEL")); logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("log level set to debug")
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
}
