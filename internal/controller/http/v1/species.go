package v1

import (
	"net/http"

	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/gin-gonic/gin"
)

type speciesRoutes struct {
	speciesUseCase usecase.Species
}

func newSpeciesRoutes(handler *gin.RouterGroup, species usecase.Species) {
	r := &speciesRoutes{species}

	h := handler.Group("/species")
	{
		h.GET("/", r.doList)
	}
}

func (r *speciesRoutes) doList(c *gin.Context) {
	title := c.Query("title")

	species, err := r.speciesUseCase.List(c.Request.Context(), title)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, schema.DataResponse{Data: species})
}
