package srv

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context, obj interface{}) bool {
	err := c.ShouldBind(obj)
	if err != nil {
		result := new(Result)
		result.Code = http.StatusBadRequest
		result.Err = err
		result.Send(c)
		return true
	}
	return false
}
