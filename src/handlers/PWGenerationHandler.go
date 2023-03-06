package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "src/databases"
	model "src/models"
	service "src/services"
)

const systemKey = "buijrfbiurh"

// Authenticate handler
func HandlePasswordGeneration(dbConnection db.DBConnection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request model.PasswordGenerationRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println(err)
		}

		var passwordEntry model.PasswordEntryTag = request.ToPasswordEntryTag()
		passwordInfo, err := dbConnection.RetrievePasswordInfo(passwordEntry)

		var password string = service.GeneratePassword(request, *passwordInfo, systemKey)

		var response = model.PasswordGenerationResponse{
			Password: password,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
