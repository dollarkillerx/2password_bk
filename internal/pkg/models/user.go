package models

type User struct {
	BasicModel
	Account             string `gorm:"type:varchar(300);uniqueIndex" json:"account"`
	PublicKey           string `gorm:"type:text" json:"public_key"`
	EncryptedPrivateKey string `gorm:"type:text" json:"encrypted_private_key"`
}
