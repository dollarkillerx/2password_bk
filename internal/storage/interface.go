package storage

import (
	"github.com/dollarkillerx/gin-template/internal/pkg/models"
	"gorm.io/gorm"
)

type Interface interface {
	DB() *gorm.DB

	GetUserByAccount(account string) (*models.User, error)
	AccountRegistry(account string, name string, password string) error
}
