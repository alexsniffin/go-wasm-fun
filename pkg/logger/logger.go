package logger

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/alexsniffin/website/pkg/models"
)

func NewLogger(cfg models.Logger, environment string) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	logger := zerolog.New(os.Stdout).Level(level).With().Timestamp().Logger()
	if environment == "localhost" {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
	}

	return logger, nil
}
