package storage

import (
	"github.com/dollarkillerx/2password/internal/pkg/models"
	"gorm.io/gorm"
)

type Interface interface {
	DB() *gorm.DB

	GetUserByAccount(account string) (*models.User, error)
	AccountRegistry(account string, publicKey string, encryptedPrivateKey string) error
}
