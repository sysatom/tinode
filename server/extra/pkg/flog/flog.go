package flog

import (
	"github.com/rs/zerolog"
	"io"
	"os"
)

var l zerolog.Logger

func init() {
	var writer []io.Writer
	// console
	console := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: zerolog.TimeFieldFormat, NoColor: true}
	writer = append(writer, console)

	multi := zerolog.MultiLevelWriter(writer...)
	l = zerolog.New(multi).With().Timestamp().Logger()
}

func Debug(format string, a ...any) {
	l.Debug().Caller(1).Msgf(format, a...)
}

func Info(format string, a ...any) {
	l.Info().Caller(1).Msgf(format, a...)
}

func Warn(format string, a ...any) {
	l.Warn().Caller(1).Msgf(format, a...)
}

func Error(err error) {
	l.Error().Caller(1).Err(err)
}

func Fatal(format string, a ...any) {
	l.Fatal().Caller(1).Msgf(format, a...)
}

func Panic(format string, a ...any) {
	l.Panic().Caller(1).Msgf(format, a...)
}
