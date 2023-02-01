package v1

import (
	"github.com/VetKA-org/vetka/internal/usecase"
	"github.com/VetKA-org/vetka/pkg/logger"
	"github.com/gin-gonic/gin"
)

type usersRoutes struct {
	log          *logger.Logger
	usersUseCase usecase.Users
}

func newUsersRoutes(handler *gin.RouterGroup, log *logger.Logger, users usecase.Users) {
	r := &usersRoutes{log, users}

	h := handler.Group("/users")
	{
		h.GET("/", r.doList)
		h.POST("/", r.doRegister)
	}
}

func (r *usersRoutes) doList(handler *gin.Context) {
}

func (r *usersRoutes) doRegister(handler *gin.Context) {
}
