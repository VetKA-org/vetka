package v1

import (
	"net/http"

	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/gin-gonic/gin"
)

type speciesRoutes struct {
	log            *logger.Logger
	speciesUseCase usecase.Species
}

func newSpeciesRoutes(handler *gin.RouterGroup, log *logger.Logger, species usecase.Species) {
	r := &speciesRoutes{log, species}

	h := handler.Group("/species")
	{
		h.GET("/", r.doList)
	}
}

func (r *speciesRoutes) doList(c *gin.Context) {
	species, err := r.speciesUseCase.List(c.Request.Context())
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, schema.DataResponse{Data: species})
}
