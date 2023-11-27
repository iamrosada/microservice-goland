// main.go
package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gin-gonic/gin"
	"github.com/iamrosada/microservice-goland/promo-service/api"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/entity"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/infra/repository"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/usecase"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Set up PostgreSQL database connection
	dbPath := "./db/main.db"
	sqlDB, err := sql.Open("sqlite3", dbPath)
	// sqlDB, err := sql.Open("postgres", "postgres://user_promo:password_promo@db_promo:5434/promo?sslmode=disable")
	// sqlDB, err := sql.Open("postgres", "postgres://user:password@db:5435/users?sslmode=disable")

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

	err = gormDB.AutoMigrate(&entity.Promotion{}, &entity.Code{})
	if err != nil {
		panic(err)
	}

	gormDB.Callback().Create().Before("gorm:before_create").Register("serializeCodes", serializeCodes)
	gormDB.Callback().Query().After("gorm:after_query").Register("deserializeCodes", deserializeCodes)
	// Create repositories and use cases
	promoRepository := repository.NewPromoRepository(gormDB)
	promoUsingUsecase := usecase.NewPromoUsecase(promoRepository)

	// Create handlers
	promoHandlers := api.NewPromoHandlers(promoUsingUsecase)

	// Set up Gin router
	router := gin.Default()

	// Set up user routes
	promoHandlers.SetupRouter(router)

	// Start the server
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
	}
}

func serializeCodes(db *gorm.DB) {
	if serializable, ok := db.Statement.Dest.(entity.Serializable); ok {
		serializable.BeforeSave()
	}
}

func deserializeCodes(db *gorm.DB) {
	if serializable, ok := db.Statement.Dest.(entity.Serializable); ok {
		serializable.AfterFind()
	}
}
