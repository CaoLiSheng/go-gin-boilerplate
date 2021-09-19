package web

import (
	srv "go-gin-boilerplate/server"
	"go-gin-boilerplate/web/test"

	"github.com/gin-gonic/gin"
)

func Setup(e *gin.Engine) {
	superAuth := srv.BearerAuth(srv.SuperUserType)
	e.POST("/login/super", superAuth.LoginHandler)
	e.POST("/refresh/super", superAuth.RefreshHandler)
	// e.POST("/logout/super", superAuth.LogoutHandler)
	super := e.Group("/super")
	super.Use(superAuth.MiddlewareFunc())
	{
		super.GET("/test", test.API)
	}
}
