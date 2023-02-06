package v1

import (
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeErrorResponse(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, errorResponse{err.Error()})
}
