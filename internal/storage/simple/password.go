package simple

import (
	"github.com/dollarkillerx/2password/internal/pkg/models"
	"github.com/rs/xid"

	"log"
)

func (s *Simple) PasswordDataInfo(account string) (pos models.PasswordDataInfo, err error) {
	err = s.db.Model(&models.PasswordOption{}).
		Where("account = ?", account).
		Where("type = ?", models.LoginType).Count(&pos.LoginTypeCount).Error
	if err != nil {
		log.Println(err)
		return
	}

	err = s.db.Model(&models.PasswordOption{}).
		Where("account = ?", account).
		Where("type = ?", models.CardType).Count(&pos.CardCount).Error
	if err != nil {
		log.Println(err)
		return
	}

	err = s.db.Model(&models.PasswordOption{}).
		Where("account = ?", account).
		Where("type = ?", models.IdentityType).Count(&pos.IdentityTypeCount).Error
	if err != nil {
		log.Println(err)
		return
	}

	err = s.db.Model(&models.PasswordOption{}).
		Where("account = ?", account).
		Where("type = ?", models.NoteType).Count(&pos.NoteTypeCount).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (s *Simple) PasswordOptionList(account string, pType models.PasswordType) (pos []models.PasswordOption, err error) {
	err = s.db.Model(&models.PasswordOption{}).
		Where("account = ?", account).
		Where("type = ?", pType).Order("created_at desc").Find(&pos).Error
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return
}

func (s *Simple) PasswordData(account string, pID string) (pos models.PasswordOption, err error) {
	err = s.db.Model(&models.PasswordOption{}).
		Where("account = ?", account).
		Where("id = ?", pID).First(&pos).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (s *Simple) DeletePasswordData(account string, pID string) (err error) {
	err = s.db.Model(&models.PasswordOption{}).
		Where("account = ?", account).
		Where("id = ?", pID).Delete(&models.PasswordOption{}).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (s *Simple) AddPasswordData(account string, pType models.PasswordType, payload string) (err error) {
	err = s.db.Model(&models.PasswordOption{}).Create(&models.PasswordOption{
		BasicModel: models.BasicModel{
			ID: xid.New().String(),
		},
		Account: account,
		Type:    pType,
		Payload: payload,
	}).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (s *Simple) UpdatePasswordData(id string, account string, payload string) (err error) {
	err = s.db.Model(&models.PasswordOption{}).
		Where("account = ?", account).
		Where("id = ?", id).Updates(map[string]interface{}{
		"payload": payload,
	}).Error
	if err != nil {
		log.Println(err)
		return
	}

	return
}
