package services

import (
	"strings"

	model "src/models"

	log "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
)

func InitConfig() (model.AppConfiguration, error) {
	viper.AddConfigPath("../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Error("Unable to read env file")
		return model.AppConfiguration{}, err
	}

	var dbConfig model.DatabaseConfiguration
	if err := viper.Unmarshal(&dbConfig); err != nil {
		log.Error("Unable to parse database config")
		return model.AppConfiguration{}, err
	}

	var secretKeys model.SecretKeys
	if err := viper.Unmarshal(&secretKeys); err != nil {
		log.Error("Unable to parse keys")
		return model.AppConfiguration{}, err
	}

	var awsConfig model.AWSConfiguration
	if err := viper.Unmarshal(&awsConfig); err != nil {
		log.Error("Unable to parse aws config")
		return model.AppConfiguration{}, err
	}

	log.Info("Successfully loaded configs")

	privkey := strings.Split(secretKeys.RSAKeyFile, "/")
	err := DownloadS3Object(
		awsConfig.S3Bucket,
		privkey[len(privkey)-1],
		secretKeys.RSAKeyFile,
		awsConfig.Region,
	)

	if err != nil {
		log.Error("Unable to download objects from S3")
		return model.AppConfiguration{}, err
	}

	return model.AppConfiguration{
		Database:  dbConfig,
		SecretKey: secretKeys,
		AWSConfig: awsConfig,
	}, nil
}
