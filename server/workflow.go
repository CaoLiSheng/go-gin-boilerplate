package srv

import (
	"go-gin-boilerplate/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Result) Send(c *gin.Context) {
	c.JSON(r.Code, r)

	if r.Err != nil {
		log.Println("error occurred:\n", r)
	}
}

func Do(c *gin.Context, opts *db.JobOptions, job Job, simple, auto bool) (result *Result) {
	opts.Job = func (core *db.Core) {
		result = job(core)
	}
	opts.Fail = func (err error) {
		result = new(Result)
		result.Code = http.StatusServiceUnavailable
		result.Err = err
	}

	if simple {
		MustGet(c).DoSimple(opts)
	} else {
		MustGet(c).Do(opts)
	}

	if auto {
		result.Send(c)
	}

	return
}
