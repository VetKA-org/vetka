// Implements v1 routing paths.
package v1

import (
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, log *logger.Logger, useCases *usecase.UseCases) {
	// Common midleware.
	handler.Use(gin.Recovery())

	// Routers
	h := handler.Group("/api/v1")
	{
		newAuthRoutes(h, log, useCases.Auth)
		newUsersRoutes(h, log, useCases.Users)
		newPatientsRoutes(h, log, useCases.Patients)
		newAppointmentsRoutes(h, log, useCases.Appointments)
		newQueueRoutes(h, log, useCases.Queue)
	}
}
