package server

import (
	"log"
	"time"

	"github.com/dollarkillerx/2password/internal/conf"
	"github.com/dollarkillerx/2password/internal/pkg/errs"
	"github.com/dollarkillerx/2password/internal/pkg/request"
	"github.com/dollarkillerx/2password/internal/pkg/response"
	"github.com/dollarkillerx/2password/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) userLogin(ctx *gin.Context) {
	var payload request.UserLogin
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	captchaOK := checkImgCaptcha(s.cache, payload.CaptchaID, payload.CaptchaCode)
	if !captchaOK {
		response.Return(ctx, errs.CaptchaCode2)
		return
	}

	uc, err := s.storage.GetUserByAccount(payload.Account)
	if err != nil {
		response.Return(ctx, errs.LoginFailed)
		return
	}

	if utils.GetPassword(uc.Password, conf.CONF.Salt) != payload.Password {
		response.Return(ctx, errs.LoginFailed)
		return
	}

	token, err := utils.JWT.CreateToken(request.AuthJWT{
		Account: uc.Account,
		Name:    uc.Name,
	}, time.Now().Add(time.Hour*24*600).Unix())
	if err != nil {
		response.Return(ctx, errs.SystemError)
		return
	}

	response.Return(ctx, response.JWT{
		JWT: token,
	})
}

func (s *Server) userRegistry(ctx *gin.Context) {
	var payload request.UserRegistry
	err := ctx.ShouldBindJSON(&payload)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	captchaOK := checkImgCaptcha(s.cache, payload.CaptchaID, payload.CaptchaCode)
	if !captchaOK {
		response.Return(ctx, errs.CaptchaCode2)
		return
	}

	err = s.storage.AccountRegistry(payload.Account, payload.Name, utils.GetPassword(payload.Password, conf.CONF.Salt))
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}
