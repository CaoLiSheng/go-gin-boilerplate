package test

import (
	"go-gin-boilerplate/db"
	srv "go-gin-boilerplate/server"

	"github.com/gin-gonic/gin"
)

// API :
func API(c *gin.Context) {
	req := new(req)
	if srv.BadRequest(c, req) {
		return
	}

	srv.Do(c, db.NewJobOpts(true, true), handler(req))
}

