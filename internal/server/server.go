package server

import (
	"github.com/afocus/captcha"
	"github.com/dollarkillerx/gin-template/internal/conf"
	"github.com/dollarkillerx/gin-template/internal/middleware"
	"github.com/dollarkillerx/gin-template/internal/storage"
	"github.com/dollarkillerx/gin-template/internal/storage/simple"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"

	"time"
)

type Server struct {
	app     *gin.Engine
	cache   *cache.Cache
	storage storage.Interface
	captcha *captcha.Captcha
}

func NewServer() *Server {
	ser := &Server{
		cache: cache.New(15*time.Minute, 30*time.Minute),
		app:   gin.New(),
	}

	ser.captchaInit()

	return ser
}

func (s *Server) Run() error {
	newSimple, err := simple.NewSimple(&conf.CONF.PgSQLConfig)
	if err != nil {
		return err
	}

	s.storage = newSimple

	s.app.Use(middleware.SetBasicInformation())
	s.app.Use(middleware.Cors())
	s.app.Use(middleware.HttpRecover())

	s.router()

	return s.app.Run(conf.CONF.ListenAddr)
}
