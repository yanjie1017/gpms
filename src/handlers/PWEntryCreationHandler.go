package handlers

import (
	"encoding/json"
	"net/http"

	db "src/databases"
	model "src/models"
	util "src/utils"

	log "github.com/sirupsen/logrus"
)

func HandlePasswordEntryCreation(dbConnection db.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contextLogger := log.WithFields(log.Fields{
			"endpoint": util.PASSWORD_ENTRY_CREATION_ENDPOINT,
		})

		w.Header().Set("Content-Type", "application/json")

		var request model.EntryCreationRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to decode request")
			w.WriteHeader(http.StatusBadRequest)
		}

		var passwordInfo model.PasswordInfo = request.ToPasswordInfo()
		entryId, err := dbConnection.CreatePasswordEntry(passwordInfo)

		if err != nil {
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error("Unable to create password entry in database")
			w.WriteHeader(http.StatusBadRequest)
		}

		var response = model.EntryCreationResponse{
			EntryId:          entryId,
			EntryReferenceId: request.EntryReferenceId,
		}

		contextLogger.WithFields(log.Fields{
			"entry_id":           response.EntryId,
			"entry_reference_id": response.EntryReferenceId,
		}).Info("Created password entry in database")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
