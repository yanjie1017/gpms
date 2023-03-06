package models

type EntryCreationResponse struct {
	EntryId          int64  `json:"entryId"`
	EntryReferenceId string `json:"entryReferenceId"`
}

type PasswordGenerationResponse struct {
	Password string `json:"password"`
}
