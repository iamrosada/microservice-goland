package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/iamrosada/microservice-goland/promo-service/api"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/entity"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/infra/repository"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/usecase"
	_ "github.com/mattn/go-sqlite3" // SQLite driver

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Set up SQLite database connection
	dbPath := "./db/main.db"
	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	_, err = os.Stat(dbPath)
	if os.IsNotExist(err) {
		err = os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			panic(err)
		}

		file, err := os.Create(dbPath)
		if err != nil {
			panic(err)
		}
		file.Close()
	}

	// Create Gorm connection
	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// AutoMigrate ensures that the database schema is up-to-date
	err = gormDB.AutoMigrate(&entity.Promotion{}, &entity.CodesPromo{})
	if err != nil {
		panic(err)
	}

	// Callback functions for serialization and deserialization
	gormDB.Callback().Create().Before("gorm:before_create").Register("serializeCodes", serializeCodes)
	gormDB.Callback().Query().After("gorm:after_query").Register("deserializeCodes", deserializeCodes)

	// Create repositories and use cases
	promoRepository := repository.NewPromoRepository(gormDB)
	promoUsingUsecase := usecase.NewPromoUsecase(promoRepository)

	// Create handlers
	promoHandlers := api.NewPromoHandlers(promoUsingUsecase)

	// Set up Gin router
	router := gin.Default()

	// Set up promo routes
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
