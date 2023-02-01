package v1

import (
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
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
}
