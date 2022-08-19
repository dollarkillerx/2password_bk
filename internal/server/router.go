package server

import (
	"github.com/gin-gonic/gin"
)

func (s *Server) router() {
	s.app.Use(gin.Logger())

	v1 := s.app.Group("/api/v1")
	v1.GET("captcha", s.showCaptcha)
	v1.POST("login", s.userLogin)
	v1.POST("registry", s.userRegistry)
}
