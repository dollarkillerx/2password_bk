package server

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/dollarkillerx/2password/internal/pkg/errs"
	"github.com/dollarkillerx/2password/internal/pkg/request"
	"github.com/dollarkillerx/2password/internal/pkg/response"
	"github.com/dollarkillerx/2password/internal/utils"
	"github.com/gin-gonic/gin"
)

func (s *Server) userInfo(ctx *gin.Context) {
	account := ctx.Param("account")

	uc, err := s.storage.GetUserByAccount(account)
	if err != nil {
		response.Return(ctx, errs.BadRequest)
		return
	}

	response.Return(ctx, gin.H{
		"public_key":            uc.PublicKey,
		"encrypted_private_key": uc.EncryptedPrivateKey,
	})
}

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
		log.Println(err)
		response.Return(ctx, errs.LoginFailed)
		return
	}

	// 验证签名
	split := strings.Split(payload.Sign, ":")
	if len(split) != 2 {
		log.Println("len(split) != 2")
		response.Return(ctx, errs.BadRequest)
		return
	}

	data, err := utils.Base64Decode(split[0])
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	sig := split[1]
	decode, err := utils.Base64Decode(sig)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	err = utils.RSASignVer(split[0], decode, uc.PublicKey)
	if err != nil {
		log.Println(err)
		log.Println(string(data))
		log.Println(string(decode))
		log.Println(uc.PublicKey)
		response.Return(ctx, errs.BadRequest)
		return
	}

	var exp request.LoginExpiration
	err = json.Unmarshal(data, &exp)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.BadRequest)
		return
	}

	if time.Now().UnixMilli() > exp.Expiration {
		log.Println(exp.Expiration)
		response.Return(ctx, errs.BadRequest)
		return
	}

	token, err := utils.JWT.CreateToken(request.AuthJWT{
		Account: uc.Account,
	}, 0)
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

	err = s.storage.AccountRegistry(payload.Account, payload.PublicKey, payload.EncryptedPrivateKey)
	if err != nil {
		log.Println(err)
		response.Return(ctx, errs.SqlSystemError)
		return
	}

	response.Return(ctx, gin.H{})
}
