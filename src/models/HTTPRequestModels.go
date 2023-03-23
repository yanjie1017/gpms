package models

type GeneralRequest struct {
	Header string `json:"header"`
	Body   string `json:"body`
}

type RequestHeader struct {
	ClientId int64 `json:"clientId"`
}

type EntryCreationRequest struct {
	ClientId         int64  `json:"clientId" validate:"required"`
	EntryReferenceId string `json:"entryReferenceId" validate:"required"`
	Length           int64  `json:"passwordLength" validate:"required"`
	SiteName         string `json:"siteName" validate:"required"`
	SiteType         string `json:"siteType" validate:"required"`
	Metadata         string `json:"metadata" validate:"required"`
	Username         string `json:"username"`
}

type PasswordGenerationRequest struct {
	ClientId        int64  `json:"clientId" validate:"required"`
	EntryId         int64  `json:"entryId" validate:"required"`
	UserInput       string `json:"userInput" validate:"required"`
	GenerationToken string `json:"token" validate:"required"`
}

func (r *EntryCreationRequest) ToPasswordInfo() PasswordInfo {
	clientInfo := ClientInfo{
		Id: r.ClientId,
	}
	entryInfo := PasswordEntry{
		ReferenceId: r.EntryReferenceId,
		Length:      r.Length,
	}
	siteInfo := SiteInfo{
		Name:     r.SiteName,
		Type_:    r.SiteType,
		Metadata: r.Metadata,
		Username: r.Username,
	}
	passwordInfo := PasswordInfo{
		Client: clientInfo,
		Entry:  entryInfo,
		Site:   siteInfo,
	}
	return passwordInfo
}

func (r *PasswordGenerationRequest) ToPasswordEntryTag() PasswordEntryTag {
	return PasswordEntryTag{
		ClientId: r.ClientId,
		EntryId:  r.EntryId,
	}
}
