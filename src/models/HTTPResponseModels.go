package models

type GeneralResponse struct {
	Header string `json:"header"`
	Body   string `json:"body`
}

type ResponseHeader struct {
	Signature string `json:"signature"`
}

type EntryCreationResponse struct {
	EntryId          int64  `json:"entryId"`
	EntryReferenceId string `json:"entryReferenceId"`
}

type PasswordGenerationResponse struct {
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json: message`
}
