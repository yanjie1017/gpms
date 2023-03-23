package models

type GeneralRequest struct {
	Header string `json:"header"`
	Body   string `json:"body`
}

type RequestHeader struct {
	ClientId int64 `json:"clientId"`
}

type EntryCreationRequest struct {
	ClientId         int64  `json:"clientId"`
	EntryReferenceId string `json:"entryReferenceId"`
	Length           int64  `json:"passwordLength"`
	SiteName         string `json:"siteName"`
	SiteType         string `json:"siteType,omitempty"`
	Metadata         string `json:"metadata"`
	Username         string `json:"username,omitempty"`
}

type PasswordGenerationRequest struct {
	ClientId        int64  `json:"clientId"`
	EntryId         int64  `json:"entryId"`
	UserInput       string `json:"userInput"`
	GenerationToken string `json:"token"`
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
