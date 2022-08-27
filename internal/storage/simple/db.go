package simple

import (
	"github.com/dollarkillerx/2password/internal/conf"
	"github.com/dollarkillerx/2password/internal/pkg/models"
	"github.com/dollarkillerx/2password/internal/utils"
	"gorm.io/gorm"

	"sync"
)

type Simple struct {
	db *gorm.DB

	inventoryMu sync.Mutex
}

func NewSimple(conf *conf.PgSQLConfig) (*Simple, error) {
	sql, err := utils.InitPgSQL(conf)
	if err != nil {
		return nil, err
	}

	sql.AutoMigrate(
		&models.User{},
		&models.PasswordOption{},
	)

	return &Simple{
		db: sql,
	}, nil
}

func (s *Simple) DB() *gorm.DB {
	return s.db
}
