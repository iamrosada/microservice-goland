package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamrosada/microservice-goland/user-service/api"
	"github.com/iamrosada/microservice-goland/user-service/internal/user/infra/repository"
	"github.com/iamrosada/microservice-goland/user-service/internal/user/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Set up PostgreSQL database connection
	sqlDB, err := sql.Open("postgres", "postgres://user:password@localhost:5432/users?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	// Create Gorm connection
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Create repositories and use cases
	userRepository := repository.NewUserRepositoryPostgres(gormDB)
	createUserUsecase := usecase.NewCreateUserUseCase(userRepository)
	listUsersUsecase := usecase.NewGetAllUsersUseCase(userRepository)
	deleteUserUsecase := usecase.NewDeleteUserUseCase(userRepository)
	getUserByIDUsecase := usecase.NewGetUserByIDUseCase(userRepository)
	updateUserUsecase := usecase.NewUpdateUserUseCase(userRepository)

	// Create handlers
	userHandlers := api.NewUserHandlers(createUserUsecase, listUsersUsecase, deleteUserUsecase, getUserByIDUsecase, updateUserUsecase)

	// Set up Gin router
	router := gin.Default()

	// Set up user routes
	userHandlers.SetupRoutes(router)

	// Start the server
	err = http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err)
	}
}
