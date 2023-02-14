package v1

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var errURLHasNoParam = errors.New("parameter not found in request URL")

// Extract ID as UUID from request URL parameters.
func getParamUUID(c *gin.Context, key string) (uuid.UUID, error) {
	rawValue := c.Param(key)
	if rawValue == "" {
		return uuid.UUID{}, fmt.Errorf("%w (%s)", errURLHasNoParam, key)
	}

	return uuid.FromString(rawValue)
}
