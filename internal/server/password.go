package server

import (
	"github.com/dollarkillerx/2password/internal/pkg/errs"
	"github.com/dollarkillerx/2password/internal/pkg/models"
	"github.com/dollarkillerx/2password/internal/pkg/request"
	"github.com/dollarkillerx/2password/internal/pkg/response"
	"github.com/dollarkillerx/2password/internal/utils"
	"github.com/gin-gonic/gin"

	"log"
)

type passwordManager struct {
	s *Server
}

func PasswordManager(s *Server) *passwordManager {
	return &passwordManager{s: s}
}

func (p *passwordManager) allInfo(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	pos, err := p.s.storage.PasswordDataInfo(model.Account)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, pos)
}

func (p *passwordManager) list(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)
	pType := ctx.Query("type")
	if pType == "" {
		response.Return(ctx, errs.BadRequest)
		return
	}
	pos, err := p.s.storage.PasswordOptionList(model.Account, models.PasswordType(pType))
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, pos)
}

func (p *passwordManager) info(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)

	pID := ctx.Param("id")

	data, err := p.s.storage.PasswordData(model.Account, pID)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, data)
}

func (p *passwordManager) add(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)
	var payload request.PassAdd
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	err = p.s.storage.AddPasswordData(model.Account, payload.Type, payload.Payload)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}

func (p *passwordManager) delete(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)
	pID := ctx.Param("id")
	err := p.s.storage.DeletePasswordData(model.Account, pID)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})

}
func (p *passwordManager) update(ctx *gin.Context) {
	model := utils.GetAuthModel(ctx)
	var payload request.PassUpdate
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	err = p.s.storage.UpdatePasswordData(payload.ID, model.Account, payload.Payload)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}
