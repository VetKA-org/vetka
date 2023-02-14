package v1

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type errorResponse struct {
	Error string `json:"error"`
}

type invalidField struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"

	case "lte":
		return "Must be less than " + fe.Param()

	case "gte":
		return "Must be greater than " + fe.Param()

	case "max":
		return "Must be shorter than " + fe.Param() + " character(s)"

	case "ltefield":
		return "Must be less than " + fe.Param()

	case "gtefield":
		return "Must be greater than " + fe.Param()

	case "oneof":
		return "Must be one of: " + fe.Param()

	default:
		return "Unknown error: " + fe.Tag() + "(" + fe.Param() + ")"
	}
}

func writeErrorResponse(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, errorResponse{err.Error()})
}

func writeBindErrorResponse(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		writeErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	out := make([]invalidField, len(ve))
	for i, fe := range ve {
		out[i] = invalidField{strings.ToLower(fe.Field()), getErrorMsg(fe)}
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
}
