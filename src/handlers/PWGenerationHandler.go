package handlers

import (
	"encoding/json"
	"net/http"

	db "src/databases"
	model "src/models"
	service "src/services"
	util "src/utils"

	log "github.com/sirupsen/logrus"
)

const systemKey = "buijrfbiurh"

// Authenticate handler
func HandlePasswordGeneration(dbConnection db.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contextLogger := log.WithFields(log.Fields{
			"endpoint": util.PASSWORD_GENERATION_ENDPOINT,
		})

		var request model.PasswordGenerationRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to decode request")
			w.WriteHeader(http.StatusBadRequest)
			// TODO: error response
		}

		var passwordEntry model.PasswordEntryTag = request.ToPasswordEntryTag()
		passwordInfo, err := dbConnection.RetrievePasswordInfo(passwordEntry)

		if err != nil {
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to retrive password info from database")
			w.WriteHeader(http.StatusBadRequest)
			// TODO: error response
		}

		password, err := service.GeneratePassword(request, *passwordInfo, systemKey)

		if err != nil {
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to generate password")
			w.WriteHeader(http.StatusBadRequest)
			// TODO: error response
		}

		var response = model.PasswordGenerationResponse{
			Password: password,
		}

		contextLogger.WithFields(log.Fields{
			"client_id": request.ClientId,
			"entry_id":  request.ClientId,
		}).Info("Generated password")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
