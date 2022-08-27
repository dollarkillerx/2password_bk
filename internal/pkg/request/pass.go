package request

import "github.com/dollarkillerx/2password/internal/pkg/models"

type PassAdd struct {
	Type    models.PasswordType `json:"type" binding:"required"`
	Payload string              `json:"payload" binding:"required"`
}

type PassUpdate struct {
	ID      string `json:"id" binding:"required"`
	Payload string `json:"payload" binding:"required"`
}
