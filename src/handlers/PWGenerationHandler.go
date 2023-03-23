package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	db "src/databases"
	model "src/models"
	service "src/services"
	util "src/utils"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

func HandlePasswordGeneration(dbConnection db.DBConnection, jsonValidator *validator.Validate, secretKeys model.SecretKeys) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contextLogger := log.WithFields(log.Fields{
			"endpoint": util.PASSWORD_GENERATION_ENDPOINT,
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

		requestBody := model.PasswordGenerationRequest{}
		err = json.Unmarshal([]byte(requestBodyString), &requestBody)
		if err != nil {
			errorMsg := util.HTTP_ERROR_RESPONSE_PARSE_REQUEST
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error(errorMsg)
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		err = jsonValidator.Struct(requestBody)
		if err != nil {
			errors := err.(validator.ValidationErrors)
			errorMsg := util.HTTP_ERROR_RESPONSE_MISSING_FIELD
			e := errors[0]
			errorMsg += fmt.Sprintf("%s is %s.", e.Field(), e.Tag())
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		contextLogger = log.WithFields(log.Fields{
			"client_id": requestBody.ClientId,
			"entry_id":  requestBody.ClientId,
		})

		var passwordEntry model.PasswordEntryTag = requestBody.ToPasswordEntryTag()
		passwordInfo, err := dbConnection.RetrievePasswordInfo(passwordEntry)

		if err != nil {
			errorMsg := util.HTTP_ERROR_RESPONSE_PASSWORD_ENTRY_NOT_FOUND
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error(errorMsg)
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		password, err := service.GeneratePassword(requestBody, *passwordInfo, secretKeys.HashKey)

		if err != nil {
			errorMsg := util.HTTP_ERROR_RESPONSE_PASSWORD_ENTRY_NOT_FOUND
			contextLogger.WithFields(log.Fields{
				"error": err,
			}).Error(errorMsg)
			w = helper.ReturnErrorResponse(w, errorMsg)
			return
		}

		contextLogger.Info("Successfully generated password")

		var responseBody = model.PasswordGenerationResponse{
			Password: password,
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
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(encryptedResponse))
	}
}
