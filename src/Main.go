package main

import (
	"net/http"

	db "src/databases"
	handler "src/handlers"
	util "src/utils"

	mux "github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const port = ":8000"

func main() {
	util.SetupLogger()
	log.Info("Initialising application...")

	// Connect to database
	dbConnection := db.ConnectDB()

	// Defer closing database connection
	defer dbConnection.DB.Close()

	// Listen to HTTP requests
	log.Info("Listening to HTTP requests...")
	router := mux.NewRouter()
	router.HandleFunc(util.HEALTHCHECK_ENDPOINT, handler.HandleHealthCheck).Methods(http.MethodPost)
	router.HandleFunc(util.PASSWORD_GENERATION_ENDPOINT, handler.HandlePasswordGeneration(dbConnection)).Methods((http.MethodPost))
	router.HandleFunc(util.HEALTHCHECK_ENDPOINT, handler.HandlePasswordEntryCreation(dbConnection)).Methods(http.MethodPost)

	err := http.ListenAndServe(port, router)
	if err != nil {
		log.WithFields(log.Fields{
			"port":  port,
			"error": err,
		}).Error("Unable to handle HTTP requests")
		panic(err)
	}
}
