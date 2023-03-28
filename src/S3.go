package main

import (
	"strings"

	service "src/services"

	log "github.com/sirupsen/logrus"
)

func main() {

	appConfig, err := service.InitConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to load configs")
		panic(err)
	}

	privkey := strings.Split(appConfig.SecretKey.RSAKeyFile, "/")
	service.DownloadS3Object(
		appConfig.AWSConfig.S3Bucket,
		privkey[len(privkey)-1],
		appConfig.SecretKey.RSAKeyFile,
	)
}
