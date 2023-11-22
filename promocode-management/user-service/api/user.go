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
			users.GET("/", p.ListUsersHandler)
			users.DELETE("/", p.DeleteUserHandler)
			users.GET("/:id", p.GetUserByIDHandler)
			users.PUT("/", p.UpdateUserHandler)
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

func (p *UserHandlers) ListUsersHandler(c *gin.Context) {
	output, err := p.ListUsersUseCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (p *UserHandlers) DeleteUserHandler(c *gin.Context) {
	var input usecase.DeleteUserInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := p.DeleteUserUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (p *UserHandlers) GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")
	input := usecase.GetUserByIDInputDto{ID: id}
	output, err := p.GetUserByIDUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}

func (p *UserHandlers) UpdateUserHandler(c *gin.Context) {
	var input usecase.UpdateUserInputDto
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := p.UpdateUserUseCase.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, output)
}
