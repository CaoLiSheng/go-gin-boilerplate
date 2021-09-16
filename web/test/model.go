package test

import (
	"go-gin-boilerplate/db"
	srv "go-gin-boilerplate/server"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handler(req *req) srv.Job {
	return func(c *db.Core) *srv.Result {
		rows, err := c.DB.QueryxContext(*c.Ctx, "show tables");

		if err != nil {
			panic(err)
		}
		
		defer rows.Close()

		data := make([]res, 0)
		for rows.Next() {
			var res res
			err := rows.StructScan(&res)
			if err != nil {
				panic(err)
			}
			data = append(data, res)
		}

		return &srv.Result{ Code: http.StatusOK, Results: gin.H{"tables": data, "request": req} }
	}
}
