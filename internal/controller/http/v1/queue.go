package v1

import (
	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

type queueRoutes struct {
	log          *logger.Logger
	queueUseCase usecase.Queue
}

func newQueueRoutes(handler *gin.RouterGroup, log *logger.Logger, queue usecase.Queue) {
	r := &queueRoutes{log, queue}

	h := handler.Group("/queue")
	h.Use(authorizedAccess(log, []string{entity.Administrator, entity.Doctor, entity.HeadDoctor}))
	{
		h.GET("/", r.doList)
		h.POST("/", r.doEnqueue)
		h.POST("/:id/up", r.doMoveUp)
		h.POST("/:id/down", r.doMoveDown)
		h.DELETE("/:id", r.doDequeue)
	}
}

func (r *queueRoutes) doList(handler *gin.Context) {
}

func (r *queueRoutes) doEnqueue(handler *gin.Context) {
}

func (r *queueRoutes) doMoveUp(handler *gin.Context) {
}

func (r *queueRoutes) doMoveDown(handler *gin.Context) {
}

func (r *queueRoutes) doDequeue(handler *gin.Context) {
}
