package types

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Error(c *gin.Context, code int, err error, msg string) {
	var res Response
	res.Msg = err.Error()
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnError(code))
}

func OK(c *gin.Context, data interface{}, msg string) {
	var res Response
	res.Data = data
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnOK())
}

func PageOK(c *gin.Context, result interface{}, count int, page int, limit int, msg string) {
	var res PageResponse
	res.Data.Details = result
	res.Data.Total = count
	res.Data.Page = page
	res.Data.Limit = limit
	if msg != "" {
		res.Msg = msg
	}
	c.JSON(http.StatusOK, res.ReturnOK())
}
