package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamrosada/microservice-goland/user-service/internal/user/usecase"
)

type UserHandlers struct {
	CreateUserUseCase  *usecase.CreateUserUseCase
	ListUsersUseCase   *usecase.GetAllUsersUseCase
	DeleteUserUseCase  *usecase.DeleteUserUseCase
	GetUserByIDUseCase *usecase.GetUserByIDUseCase
	UpdateUserUseCase  *usecase.UpdateUserUseCase
}

func NewUserHandlers(
	createUserUseCase *usecase.CreateUserUseCase,
	listUsersUseCase *usecase.GetAllUsersUseCase,
	deleteUserUseCase *usecase.DeleteUserUseCase,
	getUserByIDUseCase *usecase.GetUserByIDUseCase,
	updateUserUseCase *usecase.UpdateUserUseCase,
) *UserHandlers {
	return &UserHandlers{
		CreateUserUseCase:  createUserUseCase,
		ListUsersUseCase:   listUsersUseCase,
		DeleteUserUseCase:  deleteUserUseCase,
		GetUserByIDUseCase: getUserByIDUseCase,
		UpdateUserUseCase:  updateUserUseCase,
	}
}

func (p *UserHandlers) SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/", p.CreateUserHandler)

		}
	}
}

func (p *UserHandlers) CreateUserHandler(c *gin.Context) {
	var input usecase.CreateUserInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := p.CreateUserUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, output)
}
