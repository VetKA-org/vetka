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

type usersRoutes struct {
	log          *logger.Logger
	usersUseCase usecase.Users
}

func newUsersRoutes(handler *gin.RouterGroup, log *logger.Logger, users usecase.Users) {
	r := &usersRoutes{log, users}

	h := handler.Group("/users")
	h.Use(authorizedAccess(log, []string{entity.HeadDoctor}))
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

	c.JSON(http.StatusOK, schema.ListUsersResponse{Data: users})
}

func (r *usersRoutes) doRegister(c *gin.Context) {
	var req schema.RegisterUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		writeBindErrorResponse(c, err)

		return
	}

	err := r.usersUseCase.Register(c.Request.Context(), req.Login, req.Password, req.Roles)
	if err != nil {
		if errors.Is(err, entity.ErrUserExists) {
			writeErrorResponse(c, http.StatusConflict, entity.ErrUserExists)

			return
		}

		if errors.Is(err, entity.ErrRoleNotFound) {
			writeErrorResponse(c, http.StatusBadRequest, entity.ErrRoleNotFound)

			return
		}

		writeErrorResponse(c, http.StatusInternalServerError, err)

		return
	}

	c.Status(http.StatusOK)
}
