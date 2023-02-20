package v1

import (
	"net/http"
	"time"

	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/entity"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/VetKA-org/vetka/pkg/schema"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type appointmentsRoutes struct {
	appointmentsUseCase usecase.Appointments
}

type doCreateAppointmentRequest struct {
	PatientID    uuid.UUID `json:"patient_id" binding:"required"`
	AssigneeID   uuid.UUID `json:"assignee_id" binding:"required"`
	ScheduledFor time.Time `json:"scheduled_for" time_utc:"1" binding:"required"`
	Reason       string    `json:"reason" binding:"required,max=255"`
	Details      *string   `json:"details"`
}

type doUpdateAppointmentRequest struct {
	Status entity.ApptStatus `json:"status" binding:"required,oneof=scheduled opened closed canceled"`
}

func newAppointmentsRoutes(
	handler *gin.RouterGroup,
	log *logger.Logger,
	appointments usecase.Appointments,
) {
	r := &appointmentsRoutes{appointments}

	h := handler.Group("/appointments")
	h.Use(authorizedAccess(log, []string{entity.Administrator, entity.Doctor, entity.HeadDoctor}))
	{
		h.GET("/", r.doList)
		h.POST("/", r.doCreate)
		h.PATCH("/:id", r.doUpdate)
	}
}

func (r *appointmentsRoutes) doList(c *gin.Context) {
	var patientID *uuid.UUID

	if value, ok := c.GetQuery("patient_id"); ok {
		patientUUID, err := uuid.FromString(value)
		if err != nil {
			writeErrorResponse(c, http.StatusBadRequest, err)

			return
		}

		patientID = &patientUUID
	}

	appointments, err := r.appointmentsUseCase.List(c.Request.Context(), patientID)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, schema.DataResponse{Data: appointments})
}

func (r *appointmentsRoutes) doCreate(c *gin.Context) {
	var req doCreateAppointmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindErrorResponse(c, err)

		return
	}

	err := r.appointmentsUseCase.Create(
		c.Request.Context(),
		req.PatientID,
		req.AssigneeID,
		req.ScheduledFor,
		req.Reason,
		req.Details,
	)
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}

func (r *appointmentsRoutes) doUpdate(c *gin.Context) {
	var req doUpdateAppointmentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindErrorResponse(c, err)

		return
	}

	id, err := getParamID(c)
	if err != nil {
		writeErrorResponse(c, http.StatusBadRequest, err)

		return
	}

	if err := r.appointmentsUseCase.Update(c.Request.Context(), id, req.Status); err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}
