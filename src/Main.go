package main

import (
	"net/http"
	"reflect"
	"strings"

	db "src/databases"
	handler "src/handlers"
	service "src/services"
	util "src/utils"

	"github.com/go-playground/validator/v10"
	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const port = ":8000"

func main() {
	util.SetupLogger()
	log.Info("Initialising application...")

	// Init config
	log.Info("Loading configs...")
	appConfig, err := service.InitConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Unable to load configs")
		panic(err)
	}

	// Connect to database
	dbConnection := db.ConnectDB(appConfig.Database)

	// Defer closing database connection
	defer dbConnection.DB.Close()

	// Json validator
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Listen to HTTP requests
	log.Info("Listening to HTTP requests...")
	router := mux.NewRouter()
	router.HandleFunc(util.HEALTHCHECK_ENDPOINT, handler.HandleHealthCheck).Methods(http.MethodPost)
	router.HandleFunc(util.PASSWORD_GENERATION_ENDPOINT, handler.HandlePasswordGeneration(dbConnection, v, appConfig.SecretKey)).Methods((http.MethodPost))
	router.HandleFunc(util.PASSWORD_ENTRY_CREATION_ENDPOINT, handler.HandlePasswordEntryCreation(dbConnection, v, appConfig.SecretKey)).Methods(http.MethodPost)

	err = http.ListenAndServe(port, router)
	if err != nil {
		log.WithFields(log.Fields{
			"port":  port,
			"error": err,
		}).Error("Unable to handle HTTP requests")
		panic(err)
	}
}
