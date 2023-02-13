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
		newAppointmentsRoutes(h, log, useCases.Appointments)
		newAuthRoutes(h, log, useCases.Auth)
		newPatientsRoutes(h, log, useCases.Patients)
		newQueueRoutes(h, log, useCases.Queue)
		newRolesRoutes(h, log, useCases.Roles)
		newSpeciesRoutes(h, log, useCases.Species)
		newUsersRoutes(h, log, useCases.Users)
	}
}
