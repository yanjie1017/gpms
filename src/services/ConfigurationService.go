package services

import (
	model "src/models"

	log "github.com/sirupsen/logrus"
	viper "github.com/spf13/viper"
)

func InitConfig() (model.AppConfiguration, error) {
	viper.AddConfigPath("../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.WithFields(log.Fields{
			"eror": err,
		}).Error("Unable to read env file")
		return model.AppConfiguration{}, err
	}

	var dbConfig model.DatabaseConfiguration
	if err := viper.Unmarshal(&dbConfig); err != nil {
		log.WithFields(log.Fields{
			"eror": err,
		}).Error("Unable to parse database config")
		return model.AppConfiguration{}, err
	}

	var secretKeys model.SecretKeys
	if err := viper.Unmarshal(&secretKeys); err != nil {
		log.WithFields(log.Fields{
			"eror": err,
		}).Error("Unable to parse keys")
		return model.AppConfiguration{}, err
	}

	log.Info("Configs loaded")

	return model.AppConfiguration{
		Database:  dbConfig,
		SecretKey: secretKeys,
	}, nil
}
