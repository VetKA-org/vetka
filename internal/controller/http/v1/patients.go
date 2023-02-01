package v1

import (
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

func (r *patientsRoutes) doList(handler *gin.Context) {
}

func (r *patientsRoutes) doRegister(handler *gin.Context) {
}

func (r *patientsRoutes) doListAppointments(handler *gin.Context) {
}
