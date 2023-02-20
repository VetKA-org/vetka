package v1

import (
	"errors"
	"net/http"

	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/gin-gonic/gin"
)

type queueRoutes struct {
	queueUseCase usecase.Queue
}

func newQueueRoutes(handler *gin.RouterGroup, log *logger.Logger, queue usecase.Queue) {
	r := &queueRoutes{queue}

	h := handler.Group("/queue")
	h.Use(authorizedAccess(log, []string{entity.Administrator, entity.Doctor, entity.HeadDoctor}))
	{
		h.GET("/", r.doList)
		h.POST("/:id", r.doEnqueue)
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

	c.JSON(http.StatusOK, schema.DataResponse{Data: patients})
}

func (r *queueRoutes) doEnqueue(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil {
		writeErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := r.queueUseCase.Enqueue(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrAptExists) {
			writeErrorResponse(c, http.StatusConflict, entity.ErrAptExists)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}

func (r *queueRoutes) doMoveUp(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil {
		writeErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := r.queueUseCase.MoveUp(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrAptNotFound) {
			writeErrorResponse(c, http.StatusNotFound, entity.ErrAptNotFound)

			return
		}

		if errors.Is(err, entity.ErrAptHasMaxPos) {
			writeErrorResponse(c, http.StatusBadRequest, entity.ErrAptHasMaxPos)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}

func (r *queueRoutes) doMoveDown(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil {
		writeErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := r.queueUseCase.MoveDown(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrAptNotFound) {
			writeErrorResponse(c, http.StatusNotFound, entity.ErrAptNotFound)

			return
		}

		if errors.Is(err, entity.ErrAptHasMinPos) {
			writeErrorResponse(c, http.StatusBadRequest, entity.ErrAptHasMinPos)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}

func (r *queueRoutes) doDequeue(c *gin.Context) {
	id, err := getParamID(c)
	if err != nil {
		writeErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := r.queueUseCase.Dequeue(c.Request.Context(), id); err != nil {
		if errors.Is(err, entity.ErrAptNotFound) {
			writeErrorResponse(c, http.StatusNotFound, entity.ErrAptNotFound)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}
