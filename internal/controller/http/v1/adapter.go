package v1

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var errURLHasNoParam = errors.New("parameter not found in request URL")

// Extract ID as UUID from request URL parameters.
func getParamID(c *gin.Context) (uuid.UUID, error) {
	rawValue := c.Param("id")
	if rawValue == "" {
		return uuid.UUID{}, fmt.Errorf("%w (id)", errURLHasNoParam)
	}

	return uuid.FromString(rawValue)
}
