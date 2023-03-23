package handlers

import (
	"encoding/json"
	"net/http"

	db "src/databases"
	model "src/models"
	util "src/utils"

	log "github.com/sirupsen/logrus"
)

func HandlePasswordEntryCreation(dbConnection db.DBConnection, secretKeys model.SecretKeys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contextLogger := log.WithFields(log.Fields{
			"endpoint": util.PASSWORD_ENTRY_CREATION_ENDPOINT,
		})

		helper := HandlerHelper{
			secretKeys:   secretKeys,
			dbConnection: dbConnection,
		}

		requestBodyString, err := helper.DecryptAndDecodeRequest(r)
		if err != nil {
			errorMsg := util.HTTP_ERROR_RESPONSE_PARSE_REQUEST
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error(errorMsg)
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		requestBody := model.EntryCreationRequest{}
		err = json.Unmarshal([]byte(requestBodyString), &requestBody)
		if err != nil {
			errorMsg := util.HTTP_ERROR_RESPONSE_PARSE_REQUEST
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error(errorMsg)
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		contextLogger = contextLogger.WithFields(log.Fields{
			"entry_reference_id": requestBody.EntryReferenceId,
		})

		var passwordInfo model.PasswordInfo = requestBody.ToPasswordInfo()
		entryId, err := dbConnection.CreatePasswordEntry(passwordInfo)

		if err != nil {
			errorMsg := util.HTTP_ERROR_RESPONSE_PASSWORD_ENTRY_NOT_CREATED
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error(errorMsg)
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		contextLogger = contextLogger.WithFields(log.Fields{
			"entry_id": entryId,
		})

		contextLogger.Info("Successfully created password entry")

		var responseBody = model.EntryCreationResponse{
			EntryId:          entryId,
			EntryReferenceId: requestBody.EntryReferenceId,
		}

		responseBodyStr, err := json.Marshal(responseBody)

		if err != nil {
			errorMsg := util.HTTP_ERROR_RESPONSE_GENERATE_RESPONSE
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error(errorMsg)
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		encryptedResponse, err := helper.EncryptAndSignResponse(string(responseBodyStr), false)

		if err != nil {
			errorMsg := util.HTTP_ERROR_RESPONSE_GENERATE_RESPONSE
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error(errorMsg)
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(encryptedResponse))
	}
}
