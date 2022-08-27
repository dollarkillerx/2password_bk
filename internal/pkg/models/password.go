package models

type PasswordType string

const (
	LoginType    PasswordType = "login"
	CardType     PasswordType = "card"
	IdentityType PasswordType = "identity"
	NoteType     PasswordType = "note"
)

type PasswordOption struct {
	BasicModel
	Account string       `gorm:"type:varchar(300);index" json:"account"`
	Type    PasswordType `gorm:"type:varchar(300);index" json:"type"`
	Payload string       `gorm:"type:text" json:"payload"`
}

type PasswordDataInfo struct {
	LoginTypeCount    int64 `json:"login_type_count"`
	CardCount         int64 `json:"card_count"`
	IdentityTypeCount int64 `json:"identity_type_count"`
	NoteTypeCount     int64 `json:"note_type_count"`
}
