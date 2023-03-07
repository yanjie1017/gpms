package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func SetupLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}
