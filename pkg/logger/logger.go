package pkg

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() {
	logrus.SetOutput(os.Stdout)

	logrus.SetFormatter(&logrus.JSONFormatter{})
}
