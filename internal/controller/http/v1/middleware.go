package v1

import (
	"regexp"
	"strings"

	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

var loginRouteRe = regexp.MustCompile(`^/api/v\d/login`)

func authentication(log *logger.Logger, secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		matches := loginRouteRe.FindStringSubmatch(c.Request.URL.Path)
		if len(matches) > 0 {
			c.Next()

			return
		}

		rawToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		rawToken = strings.TrimSpace(rawToken)

		_, err := jwt.Parse(rawToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			log.Error().Str("path", c.Request.URL.Path).Err(err).Msg("Authentication failed")
			writeUnauthorizedErrorResponse(c)

			return
		}

		c.Next()
	}
}
