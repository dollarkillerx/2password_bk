package simple

import (
	"github.com/dollarkillerx/2password/internal/pkg/models"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

func (s *Simple) GetUserByAccount(account string) (*models.User, error) {
	var uc models.User
	err := s.DB().Model(&models.User{}).
		Where("account = ?", account).First(&uc).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &uc, nil
}

func (s *Simple) AccountRegistry(account string, name string, password string) error {
	err := s.DB().Model(&models.User{}).Create(&models.User{
		BasicModel: models.BasicModel{ID: xid.New().String()},
		Account:    account,
		Name:       name,
		Password:   password,
	}).Error
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
