package v1

import (
	"errors"
	"net/http"

	"github.com/VetKA-org/vetka/internal/entity"
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

type usersRoutes struct {
	log          *logger.Logger
	usersUseCase usecase.Users
}

type doRegisterRequest struct {
	Login    string `json:"login" binding:"required,lte=128"`
	Password string `json:"password" binding:"required"`
}

func newUsersRoutes(handler *gin.RouterGroup, log *logger.Logger, users usecase.Users) {
	r := &usersRoutes{log, users}

	h := handler.Group("/users")
	{
		h.GET("/", r.doList)
		h.POST("/", r.doRegister)
	}
}

func (r *usersRoutes) doList(c *gin.Context) {
	users, err := r.usersUseCase.List(c.Request.Context())
	if err != nil {
		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, dataResponse{Data: users})
}

func (r *usersRoutes) doRegister(c *gin.Context) {
	var req doRegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindErrorResponse(c, err)

		return
	}

	if err := r.usersUseCase.Register(c.Request.Context(), req.Login, req.Password); err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			writeErrorResponse(c, http.StatusConflict, err)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}
