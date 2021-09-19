package srv

import (
	"go-gin-boilerplate/db"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r *Result) Send(c *gin.Context) {
	if r.Err != nil {
		log.Println("error occurred:\n", r)
		r.Message = r.Err.Error()
	}

	c.JSON(r.Code, r)
}

func Do(c *gin.Context, opts *db.JobOptions, job Job) (result *Result) {
	Job := func (core *db.Core) {
		result = job(core)
	}
	Fail := func (err error) {
		result = new(Result)
		result.Code = http.StatusServiceUnavailable
		result.Err = err
	}

	if opts.Simple {
		MustGet(c).DoSimple(opts, Job, Fail)
	} else {
		MustGet(c).Do(opts, Job, Fail)
	}

	if opts.Auto {
		result.Send(c)
	}

	return
}
