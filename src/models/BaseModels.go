package models

type ClientInfo struct {
	Id int64
	// name            string
	// contact_number  string
	// email_address   string
	// mailing_address string
	// postal_code     string
	// country         string
	// is_active       bool
	// created_at      time.Time
	// updated_at      time.Time
}

type PasswordEntry struct {
	// Id          int64
	ReferenceId string
	// ClientId    int64
	// Site_id     int64
	Length int64
	// created_at       time.Time
	// updated_at       time.Time
	// last_accessed_at time.Time
}

type SiteInfo struct {
	// Id       int64
	Name     string
	Type_    string
	Metadata string
	Username string
	// created_at   time.Time
	// updated_at   time.Time
}
