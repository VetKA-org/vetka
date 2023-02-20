package v1

import (
	"net/http"

	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/gin-gonic/gin"
)

type rolesRoutes struct {
	rolesUseCase usecase.Roles
}

func newRolesRoutes(handler *gin.RouterGroup, roles usecase.Roles) {
	r := &rolesRoutes{roles}

	h := handler.Group("/roles")
	{
		h.GET("/", r.doList)
	}
}

func (r *rolesRoutes) doList(c *gin.Context) {
	name := c.Query("name")

	roles, err := r.rolesUseCase.List(c.Request.Context(), name)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, schema.ListRolesResponse{Data: roles})
}
