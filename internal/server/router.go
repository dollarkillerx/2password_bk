package server

import (
	"github.com/dollarkillerx/2password/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (s *Server) router() {
	s.app.Use(gin.Logger())

	v1 := s.app.Group("/api/v1")
	v1.GET("captcha", s.showCaptcha)
	v1.GET("user_info/:account", s.userInfo)
	v1.POST("login", s.userLogin)
	v1.POST("registry", s.userRegistry)

	internal := v1.Group("/internal", middleware.UAAuthorization())
	{
		pass := internal.Group("/password")
		{
			pass.GET("/all_info", PasswordManager(s).allInfo)
			pass.GET("/info/:id", PasswordManager(s).info)
			pass.POST("/add", PasswordManager(s).add)
			pass.POST("/delete/:id", PasswordManager(s).delete)
			pass.POST("/update", PasswordManager(s).update)
		}
	}
}
