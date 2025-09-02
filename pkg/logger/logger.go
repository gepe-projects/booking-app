package logger

import (
	"os"
	"time"

	"booking/pkg/config"

	"github.com/rs/zerolog"
)

type Logger struct {
	instance zerolog.Logger
}

func New(service string, config *config.Config) Logger {
	if config.App.Env == "dev" {
		return consoleLogger(service)
	} else {
		return jsonLogger(service)
	}
}

func consoleLogger(service string) Logger {
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		FormatTimestamp: func(i interface{}) string {
			if t, ok := i.(string); ok {
				return t
			}
			return time.Now().Format(time.RFC3339)
		},
	}

	l := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Str("service", service).
		Logger()

	return Logger{instance: l}
}

func jsonLogger(service string) Logger {
	l := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Caller().
		Str("service", service).
		Logger()

	return Logger{instance: l}
}

// Wrapper biar gampang akses
func (l Logger) Info(msg string) {
	l.instance.Info().CallerSkipFrame(1).Msg(msg)
}

func (l Logger) Infof(format string, v ...interface{}) {
	l.instance.Info().CallerSkipFrame(1).Msgf(format, v...)
}

func (l Logger) Warn(msg string) {
	l.instance.Warn().CallerSkipFrame(1).Msg(msg)
}

func (l Logger) Warnf(format string, v ...interface{}) {
	l.instance.Warn().CallerSkipFrame(1).Msgf(format, v...)
}

func (l Logger) Error(err error, msg string) {
	l.instance.Error().CallerSkipFrame(1).Err(err).Msg(msg)
}

func (l Logger) Errorf(err error, format string, v ...interface{}) {
	l.instance.Error().CallerSkipFrame(1).Err(err).Msgf(format, v...)
}

func (l Logger) Fatal(err error, msg string) {
	l.instance.Fatal().CallerSkipFrame(1).Err(err).Msg(msg)
}

func (l Logger) Fatalf(err error, format string, v ...interface{}) {
	l.instance.Fatal().CallerSkipFrame(1).Err(err).Msgf(format, v...)
}

func (l Logger) Debug(msg string) {
	l.instance.Debug().CallerSkipFrame(1).Msg(msg)
}

func (l Logger) Debugf(format string, v ...interface{}) {
	l.instance.Debug().CallerSkipFrame(1).Msgf(format, v...)
}
