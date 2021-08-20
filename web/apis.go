package web

import (
	"go-gin-boilerplate/web/test"

	"github.com/gin-gonic/gin"
)

func Setup(e *gin.Engine) {
	e.GET("/test", test.API)
}
