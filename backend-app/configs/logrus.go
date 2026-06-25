package configs

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

func NewLogger() *logrus.Logger {
	log := logrus.New()

	level := os.Getenv("LOG_LEVEL")
	number, _ := strconv.Atoi(level)
	log.SetLevel(logrus.Level(number))
	log.SetFormatter(&logrus.TextFormatter{})

	return log
}
