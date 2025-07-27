package logger

import (
	"github.com/rs/zerolog"
)

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
}
