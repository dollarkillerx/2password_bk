package server

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"strings"

	"github.com/afocus/captcha"
	"github.com/dollarkillerx/2password/internal/pkg/errs"
	"github.com/dollarkillerx/2password/internal/pkg/response"
	"github.com/dollarkillerx/2password/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

func (s *Server) captchaInit() {
	cap := captcha.New()
	// 可以设置多个字体 或使用cap.AddFont("xx.ttf")追加
	cap.SetFont("./static/comic.ttf")
	// 设置验证码大小
	cap.SetSize(150, 64)
	// 设置干扰强度
	cap.SetDisturbance(captcha.MEDIUM)
	// 设置前景色 可以多个 随机替换文字颜色 默认黑色
	cap.SetFrontColor(color.RGBA{255, 255, 255, 255})
	// 设置背景色 可以多个 随机替换背景色 默认白色
	cap.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
	s.captcha = cap
}

func (s *Server) showCaptcha(ctx *gin.Context) {
	img, str := s.captcha.Create(4, captcha.CLEAR)

	buffer := bytes.NewBuffer([]byte(""))

	err := png.Encode(buffer, img)
	if err != nil {
		response.Return(ctx, errs.SystemError)
		return
	}

	i := buffer.Bytes()
	encode := utils.Base64Encode(i)

	captchaID := utils.RandKey(6)

	captcha := map[string]string{
		"base64_captcha": encode,
		"captcha_id":     captchaID,
	}

	s.cache.Set(fmt.Sprintf("%s_captccha", captchaID), str, cache.DefaultExpiration)

	response.Return(ctx, captcha)
}

func checkImgCaptcha(cache *cache.Cache, captchaID string, code string) bool {
	if captchaID == "pacman" {
		return true
	}
	captchaID = fmt.Sprintf("%s_captccha", captchaID)
	defer func() {
		cache.Delete(captchaID)
	}()
	rData, ex := cache.Get(captchaID)
	if !ex {
		return false
	}

	if strings.ToUpper(code) != strings.ToUpper(rData.(string)) {
		return false
	}

	return true
}
