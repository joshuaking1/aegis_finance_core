package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors" // For better error stack formatting
)

// Init initializes the global zerolog logger.
// In a real app, you'd make 'level' configurable.
func Init(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs // Or time.RFC3339Nano for human-readable
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack // Better stack traces

	logLevel := zerolog.InfoLevel // Default
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	case "fatal":
		logLevel = zerolog.FatalLevel
	case "panic":
		logLevel = zerolog.PanicLevel
	}
	zerolog.SetGlobalLevel(logLevel)

	// For development, pretty console output is nice.
	// For production, you'd typically remove this or make it configurable
	// and just output JSON.
	// if os.Getenv("APP_ENV") == "development" {
	// For now, let's keep it pretty for easier visual inspection
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339Nano})
	// } else {
	// log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger() // Default JSON output
	// }

	log.Info().Str("log_level", logLevel.String()).Msg("Logger initialized")
}