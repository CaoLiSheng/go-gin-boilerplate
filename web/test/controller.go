package test

import (
	"go-gin-boilerplate/db"
	srv "go-gin-boilerplate/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

// API :
func API(c *gin.Context) {
	srv.Do(c, db.NewJobOpts(), func(c *db.Core) *srv.Result {
		return &srv.Result{ Code: http.StatusOK }
	}, true, true)
}

