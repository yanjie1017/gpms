package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "src/databases"
	models "src/models"
)

// Authenticate handler
func HandlePasswordEntryCreation(dbConnection db.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.EntryCreationRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println(err)
		}

		var passwordInfo models.PasswordInfo = request.ToPasswordInfo()
		entryId, err := dbConnection.CreatePasswordEntry(passwordInfo)

		var response = models.EntryCreationResponse{
			EntryId:          entryId,
			EntryReferenceId: request.EntryReferenceId,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}
}
