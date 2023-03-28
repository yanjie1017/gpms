package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	db "src/databases"
	model "src/models"
	service "src/services"

	log "github.com/sirupsen/logrus"
)

type HandlerHelper struct {
	dbConnection      db.DBConnection
	secretKeys        model.SecretKeys
	clientCredentials model.ClientAuthentication
}

func (h *HandlerHelper) DecryptAndDecodeRequest(r *http.Request) (string, error) {
	requestBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("Unable to read request")
		return "", err
	}
	requestString := string(requestBytes)
	requestString = strings.TrimPrefix(requestString, "\"")
	requestString = strings.TrimSuffix(requestString, "\"")

	decryptedRequest, err := service.DecryptPayload(requestString, h.secretKeys.SharedSecret)
	if err != nil {
		log.Error("Unable to decrypt request")
		return "", err
	}

	request := model.GeneralRequest{}
	err = json.Unmarshal([]byte(decryptedRequest), &request)
	if err != nil {
		log.Error("Unable to decode request")
		return "", err
	}

	requestHeader := model.RequestHeader{}
	err = json.Unmarshal([]byte(request.Header), &requestHeader)
	if err != nil {
		log.Error("Unable to decode request header")
		return "", err
	}

	encryptedApiKey, err := h.dbConnection.RetrieveAPIKey(requestHeader.ClientId)
	if err != nil {
		log.Error("Unable to retrieve api key")
		return "", err
	}

	apiKey, err := service.DecryptPayload(encryptedApiKey.APIKey, h.secretKeys.SharedSecret)
	if err != nil {
		log.Error("Unable to decrypt api key")
		return "", err
	}

	requestBody, err := service.DecryptPayload(request.Body, apiKey)
	if err != nil {
		log.Error("Unable to decrypt request body")
		return "", err
	}

	h.clientCredentials = model.ClientAuthentication{
		ClientId: requestHeader.ClientId,
		APIKey:   apiKey,
	}

	return requestBody, nil
}

func (h *HandlerHelper) EncryptAndSignResponse(body string, isError bool) (string, error) {
	contextLogger := log.WithFields(log.Fields{
		"isError": isError,
	})

	signature, err := service.SignRSASHA(h.secretKeys.SignatureMsg, h.secretKeys.RSAKeyFile)
	if err != nil {
		contextLogger.Error("Unable to sign response")
		return "", err
	}

	header := model.ResponseHeader{
		Signature: signature,
	}

	headerStr, err := json.Marshal(header)
	if err != nil {
		contextLogger.Error("Unable to encode response header")
		return "", err
	}

	response := model.GeneralResponse{
		Header: string(headerStr),
		Body:   body,
	}

	responseStr, err := json.Marshal(response)
	if err != nil {
		contextLogger.Error("Unable to encode response")
		return "", err
	}

	var encryptedResponse string
	if isError {
		encryptedResponse, err = service.EncryptPayload(string(responseStr), h.secretKeys.SharedSecret)
	} else {
		encryptedResponse, err = service.EncryptPayload(string(responseStr), h.clientCredentials.APIKey)
	}

	if err != nil {
		contextLogger.Error("Unable to encrypt response")
		return "", err
	}

	return encryptedResponse, nil
}

func (h *HandlerHelper) ReturnErrorResponse(w http.ResponseWriter, body string) http.ResponseWriter {
	errorBody := model.ErrorResponse{
		Message: body,
	}

	errorStr, err := json.Marshal(errorBody)
	if err != nil {
		log.Error("Unable to encode error message")
	}

	encryptedResponse, err := h.EncryptAndSignResponse(string(errorStr), true)
	if err != nil {
		log.Error("Unable to encrypt error response")
		encryptedResponse = ""
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(encryptedResponse))
	return w
}
