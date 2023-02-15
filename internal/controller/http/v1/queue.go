package v1

import (
	"net/http"

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

func (r *queueRoutes) doList(c *gin.Context) {
	patients, err := r.queueUseCase.List(c.Request.Context())
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, dataResponse{Data: patients})
}

func (r *queueRoutes) doEnqueue(c *gin.Context) {
}

func (r *queueRoutes) doMoveUp(c *gin.Context) {
}

func (r *queueRoutes) doMoveDown(c *gin.Context) {
}

func (r *queueRoutes) doDequeue(c *gin.Context) {
}
