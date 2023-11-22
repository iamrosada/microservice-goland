package repository

import (
	"github.com/iamrosada/microservice-goland/user-service/internal/user/entity"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserRepositoryPostgres struct {
	DB *gorm.DB
}

func NewUserRepositoryPostgres(db *gorm.DB) *UserRepositoryPostgres {
	return &UserRepositoryPostgres{DB: db}
}

func (r *UserRepositoryPostgres) Create(user *entity.User) error {
	return r.DB.Create(user).Error
}
