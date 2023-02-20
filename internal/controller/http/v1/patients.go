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

type patientsRoutes struct {
	patientsUseCase usecase.Patients
}

func newPatientsRoutes(handler *gin.RouterGroup, log *logger.Logger, patients usecase.Patients) {
	r := &patientsRoutes{patients}

	h := handler.Group("/patients")
	h.Use(authorizedAccess(log, []string{entity.Administrator, entity.Doctor, entity.HeadDoctor}))
	{
		h.GET("/", r.doList)
		h.POST("/", r.doRegister)
	}
}

func (r *patientsRoutes) doList(c *gin.Context) {
	patients, err := r.patientsUseCase.List(c.Request.Context())
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, schema.ListPatientsResponse{Data: patients})
}

func (r *patientsRoutes) doRegister(c *gin.Context) {
	var req schema.RegisterPatientRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindErrorResponse(c, err)

		return
	}

	err := r.patientsUseCase.Register(
		c.Request.Context(),
		req.Name,
		req.SpeciesID,
		req.Gender,
		req.Breed,
		req.Birth,
		req.Aggressive,
		req.VaccinatedAt,
		req.SterilizedAt,
	)
	if err != nil {
		if errors.Is(err, entity.ErrPatientExists) {
			writeErrorResponse(c, http.StatusConflict, entity.ErrPatientExists)

			return
		}

		if errors.Is(err, entity.ErrSpeciesNotFound) {
			writeErrorResponse(c, http.StatusBadRequest, entity.ErrSpeciesNotFound)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}
