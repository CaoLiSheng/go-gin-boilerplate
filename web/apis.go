package web

import (
	srv "go-gin-boilerplate/server"
	"go-gin-boilerplate/web/test"

	"github.com/gin-gonic/gin"
)

func Setup(e *gin.Engine) {
	superAuth := srv.BearerAuth(srv.SuperUserType)
	e.GET("/login/super", superAuth.LoginHandler)
	e.GET("/logout/super", superAuth.LogoutHandler)
	e.GET("/refresh/super", superAuth.RefreshHandler)
	super := e.Group("/super")
	super.Use(superAuth.MiddlewareFunc())
	{
		super.GET("/test", test.API)
	}
}
