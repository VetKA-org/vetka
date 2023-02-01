package v1

import (
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

type appointmentsRoutes struct {
	log                 *logger.Logger
	appointmentsUseCase usecase.Appointments
}

func newAppointmentsRoutes(handler *gin.RouterGroup, log *logger.Logger, appointments usecase.Appointments) {
	r := &appointmentsRoutes{log, appointments}

	h := handler.Group("/appointments")
	{
		h.GET("/", r.doList)
		h.POST("/", r.doCreate)
		h.PATCH("/:id", r.doUpdate)
	}
}

func (r *appointmentsRoutes) doList(handler *gin.Context) {
}

func (r *appointmentsRoutes) doCreate(handler *gin.Context) {
}

func (r *appointmentsRoutes) doUpdate(handler *gin.Context) {
}
