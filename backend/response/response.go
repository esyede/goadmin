package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Response(c *gin.Context, httpStatus int, code int, data gin.H, message string) {
	c.JSON(httpStatus, gin.H{"code": code, "data": data, "message": message})
}

func Success(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusOK, 200, data, message)
}

func Fail(c *gin.Context, data gin.H, message string) {
	Response(c, http.StatusBadRequest, 400, data, message)
}
