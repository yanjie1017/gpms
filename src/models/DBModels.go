package models

type PasswordInfo struct {
	Client ClientInfo
	Entry  PasswordEntry
	Site   SiteInfo
}

type PasswordEntryTag struct {
	ClientId int64
	EntryId  int64
}

type ClientAuthentication struct {
	ClientId int64
	APIKey   string //encrypted
}

type PasswordEntryCreationResult struct {
	EntryID int64
	Error   error
}

type PasswordGenerationInfo struct {
	Length   int64
	Metadata string
}
