package v1

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/pkg/compression"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

var anonymousRouteRe = regexp.MustCompile(`^/api/v\d/login`)

const _rolesKey = "roles"

func authenticatedAccess(log *logger.Logger, secret entity.Secret) gin.HandlerFunc {
	return func(c *gin.Context) {
		matches := anonymousRouteRe.FindStringSubmatch(c.Request.URL.Path)
		if len(matches) > 0 {
			c.Next()

			return
		}

		rawToken := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		rawToken = strings.TrimSpace(rawToken)

		if rawToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse{"unauthorized"})

			return
		}

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

func DecompressRequest(c *gin.Context) {
	encoding := c.GetHeader("Content-Encoding")
	if encoding == "" {
		c.Next()

		return
	}

	decoder, err := compression.NewDecoder(c.Request.Body, encoding)
	if err != nil {
		if errors.Is(err, compression.ErrEncodingNotSupported) {
			writeErrorResponse(c, http.StatusNotAcceptable, compression.ErrEncodingNotSupported)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	defer decoder.Close()
	c.Request.Body = decoder

	c.Next()
}

func CompressResponse(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		encoding := c.GetHeader("Accept-Encoding")
		if encoding == "" {
			c.Next()

			return
		}

		encoder, err := compression.NewEncoder(log, c.Writer, encoding)
		if err != nil {
			writeErrorResponse(c, http.StatusInternalServerError, err)

			return
		}

		c.Writer = encoder

		defer func() {
			encoder.Close()
			c.Header("Content-Length", strconv.Itoa(c.Writer.Size()))
		}()

		c.Next()
	}
}
