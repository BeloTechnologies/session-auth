package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	log     *logrus.Logger
	logOnce sync.Once
)

// InitLogger initializes the logger and ensures it's a singleton
func InitLogger() *logrus.Logger {
	logOnce.Do(func() {
		log = logrus.New()

		log.SetLevel(logrus.InfoLevel)

		log.SetOutput(os.Stdout)
		//file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		//if err == nil {
		//	log.Out = file
		//} else {
		//	log.Info("Failed to log to file, using default stderr")
		//}

		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})

		// For JSON format:
		// log.SetFormatter(&logrus.JSONFormatter{})
	})

	return log
}
