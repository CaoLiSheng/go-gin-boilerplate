package srv

import (
	"errors"
	"go-gin-boilerplate/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

const middlewareKey = "__db-core__"

// Middleware is a gin middleware about db context
func Middleware(core *db.Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(middlewareKey, core)
	}
}

// MustGet is a must getter for Core from Middleware
func MustGet(c *gin.Context) *db.Core {
	defer func() {
		if err := recover(); err != nil {
			r := &Result{
				Code: http.StatusInternalServerError,
				Err:  errors.New("无法获取数据库操作上下文"),
			}
			r.Send(c)
		}
	}()
	return c.MustGet(middlewareKey).(*db.Core)
}
