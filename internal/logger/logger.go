package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Infof(format string, v ...any) {
	log.Info().Msgf(format, v...)
}

func Error(msg string) {
	log.Error().Msg(msg)
}

func Err(err error, msg string) {
	log.Err(err).Msg(msg)
}

func Errorf(format string, v ...any) {
	log.Error().Msgf(format, v...)
}

func Warn(msg string) {
	log.Warn().Msg(msg)
}

func Warnf(format string, v ...any) {
	log.Warn().Msgf(format, v...)
}
