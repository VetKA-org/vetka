// Implements v1 routing paths.
package v1

import (
	"github.com/VetKA-org/vetka/internal/config"
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	handler *gin.Engine,
	log *logger.Logger,
	cfg *config.Config,
	useCases *usecase.UseCases,
) {
	// Common midleware.
	handler.Use(gin.Recovery())
	handler.Use(RequestsLogger(log))
	handler.Use(authenticatedAccess(log, cfg.Secret))
	handler.Use(DecompressRequest)
	handler.Use(CompressResponse(log))

	// Routers
	h := handler.Group("/api/v1")
	{
		newAppointmentsRoutes(h, log, useCases.Appointments)
		newAuthRoutes(h, useCases.Auth)
		newPatientsRoutes(h, log, useCases.Patients)
		newQueueRoutes(h, log, useCases.Queue)
		newRolesRoutes(h, useCases.Roles)
		newSpeciesRoutes(h, useCases.Species)
		newUsersRoutes(h, log, useCases.Users)
	}
}
