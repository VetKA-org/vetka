package v1

import (
	"errors"
	"net/http"
	"time"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type patientsRoutes struct {
	log             *logger.Logger
	patientsUseCase usecase.Patients
}

type doRegisterPatientRequest struct {
	Name       string        `json:"name" binding:"required,max=32"`
	SpeciesID  uuid.UUID     `json:"species_id" binding:"required"`
	Gender     entity.Gender `json:"gender" binding:"required,oneof=male female"`
	Breed      string        `json:"breed" binding:"required,max=64"`
	Birth      time.Time     `json:"birth" time_format:"2006-01-02" time_utc:"1" binding:"required"`
	Aggressive bool          `json:"aggressive"`
	// NB (alkurbatov): Ask lll to ignore struct tags, see:
	// https://github.com/walle/lll/issues/11
	VaccinatedAt *time.Time `json:"vaccinated_at" time_format:"2006-01-02" time_utc:"1" binding:"omitempty,gtefield=Birth"` //nolint:lll // tags
	SterilizedAt *time.Time `json:"sterilized_at" time_format:"2006-01-02" time_utc:"1" binding:"omitempty,gtefield=Birth"` //nolint:lll // tags
}

func newPatientsRoutes(handler *gin.RouterGroup, log *logger.Logger, patients usecase.Patients) {
	r := &patientsRoutes{log, patients}

	h := handler.Group("/patients")
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

	c.JSON(http.StatusOK, dataResponse{Data: patients})
}

func (r *patientsRoutes) doRegister(c *gin.Context) {
	var req doRegisterPatientRequest

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

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}
