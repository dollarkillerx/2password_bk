package server

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/afocus/captcha"
	"github.com/dollarkillerx/2password/internal/conf"
	"github.com/dollarkillerx/2password/internal/middleware"
	"github.com/dollarkillerx/2password/internal/storage"
	"github.com/dollarkillerx/2password/internal/storage/simple"
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

type Server struct {
	app     *gin.Engine
	cache   *cache.Cache
	storage storage.Interface
	captcha *captcha.Captcha
	cos     *cos.Client
}

func NewServer() *Server {
	u, _ := url.Parse(fmt.Sprintf("https://%s.cos.%s.myqcloud.com", conf.CONF.OSSConf.Bucket, conf.CONF.OSSConf.Region))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		//设置超时时间
		Timeout: 100 * time.Second,
		Transport: &cos.AuthorizationTransport{
			//如实填写账号和密钥，也可以设置为环境变量
			SecretID:  conf.CONF.OSSConf.SecretID,
			SecretKey: conf.CONF.OSSConf.SecretKey,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    false,
				ResponseHeader: true,
				ResponseBody:   false,
			},
		},
	})

	ser := &Server{
		cache: cache.New(15*time.Minute, 30*time.Minute),
		app:   gin.New(),
		cos:   c,
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
