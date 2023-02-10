package v1

import (
	"net/http"

	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

type patientsRoutes struct {
	log             *logger.Logger
	patientsUseCase usecase.Patients
}

func newPatientsRoutes(handler *gin.RouterGroup, log *logger.Logger, patients usecase.Patients) {
	r := &patientsRoutes{log, patients}

	h := handler.Group("/patients")
	{
		h.GET("/", r.doList)
		h.POST("/", r.doRegister)
		h.GET("/:id/appointments", r.doListAppointments)
	}
}

func (r *patientsRoutes) doList(c *gin.Context) {
	patients, err := r.patientsUseCase.List(c.Request.Context())
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, dataResponse{Data: patients})
}

func (r *patientsRoutes) doRegister(c *gin.Context) {
}

func (r *patientsRoutes) doListAppointments(c *gin.Context) {
}
