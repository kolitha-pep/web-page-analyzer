package logger

import (
	"fmt"
	"log"
	"os"

	"github.com/kolitha-pep/web-page-analyzer/internal/pkg/utils"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init() {
	fmt.Println("logger initializing...")

	logPath := os.Getenv("LOG_PATH")

	err := utils.CreateFileIfNotExists(logPath)
	if err != nil {
		panic(err)
	}
	//
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	//log.SetOutput(logFile)
	//log.SetFormatter(&log.JSONFormatter{})
	//log.SetLevel(log.FatalLevel)

	Log = logrus.New()
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetOutput(logFile)
	Log.SetLevel(logrus.InfoLevel)
}
