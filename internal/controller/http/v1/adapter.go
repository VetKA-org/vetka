package v1

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

var errURLHasNoKey = errors.New("key not found in request URL")

// Extract ID as UUID from request URL.
func paramToUUID(c *gin.Context, key string) (uuid.UUID, error) {
	rawValue := c.Param(key)
	if rawValue == "" {
		return uuid.UUID{}, fmt.Errorf("%w (%s)", errURLHasNoKey, key)
	}

	return uuid.FromString(rawValue)
}
