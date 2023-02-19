package v1

import (
	"errors"
	"net/http"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/gin-gonic/gin"
)

type authRoutes struct {
	log         *logger.Logger
	authUseCase usecase.Auth
}

func newAuthRoutes(handler *gin.RouterGroup, log *logger.Logger, auth usecase.Auth) {
	r := &authRoutes{log, auth}

	h := handler.Group("/")
	{
		h.POST("/login", r.doLogin)
	}
}

func (r *authRoutes) doLogin(c *gin.Context) {
	var req schema.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindErrorResponse(c, err)

		return
	}

	token, err := r.authUseCase.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			writeErrorResponse(c, http.StatusUnauthorized, entity.ErrInvalidCredentials)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Header("Authorization", "Bearer "+string(token))
	c.Status(http.StatusOK)
}
