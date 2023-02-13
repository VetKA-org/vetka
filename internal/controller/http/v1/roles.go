package v1

import (
	"net/http"

	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

type rolesRoutes struct {
	log          *logger.Logger
	rolesUseCase usecase.Roles
}

func newRolesRoutes(handler *gin.RouterGroup, log *logger.Logger, roles usecase.Roles) {
	r := &rolesRoutes{log, roles}

	h := handler.Group("/roles")
	{
		h.GET("/", r.doList)
	}
}

func (r *rolesRoutes) doList(c *gin.Context) {
	roles, err := r.rolesUseCase.List(c.Request.Context())
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, dataResponse{Data: roles})
}
