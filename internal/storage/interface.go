package storage

import (
	"github.com/dollarkillerx/2password/internal/pkg/models"
	"gorm.io/gorm"
)

type Interface interface {
	DB() *gorm.DB

	GetUserByAccount(account string) (*models.User, error)
	AccountRegistry(account string, publicKey string, encryptedPrivateKey string) error

	PasswordDataInfo(account string) (pos models.PasswordDataInfo, err error)
	PasswordOptionList(account string, pType models.PasswordType) (pos []models.PasswordOption, err error)

	PasswordData(account string, pID string) (pos models.PasswordOption, err error)
	DeletePasswordData(account string, pID string) (err error)
	AddPasswordData(account string, pType models.PasswordType, payload string) (err error)
	UpdatePasswordData(id string, account string, payload string) (err error)
}
