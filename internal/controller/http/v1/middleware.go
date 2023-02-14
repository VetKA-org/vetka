package v1

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

var anonymousRouteRe = regexp.MustCompile(`^/api/v\d/login`)

const _rolesKey = "roles"

func authenticatedAccess(log *logger.Logger, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		matches := anonymousRouteRe.FindStringSubmatch(c.Request.URL.Path)
		if len(matches) > 0 {
			c.Next()

			return
		}

		rawToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		rawToken = strings.TrimSpace(rawToken)

		decodedToken, err := entity.DecodeToken(rawToken, secret)
		if err != nil {
			log.Error().
				Str("path", c.Request.URL.Path).
				Err(err).
				Msg("authenticatedAccess - entity.DecodeToken")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"unauthorized"})

			return
		}

		c.Set(_rolesKey, decodedToken.Roles)

		c.Next()
	}
}

func authorizedAccess(log *logger.Logger, allowedRoles []string) gin.HandlerFunc {
	if len(allowedRoles) == 0 {
		log.Panic().Msg("authorization was requested but no roles provided")
	}

	allowed := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(c *gin.Context) {
		for _, role := range c.GetStringSlice(_rolesKey) {
			if _, ok := allowed[role]; ok {
				c.Next()

				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, errorResponse{"forbidden"})
	}
}
