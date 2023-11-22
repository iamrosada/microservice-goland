package entity

import "github.com/google/uuid"

type UserRepository interface {
	Create(user *User) error
	FindAll() ([]*User, error)
	Update(user *User) error
	DeleteByID(id string) error
	GetByID(id string) (*User, error)
}

type User struct {
	ID    string
	Name  string
	Email string
}

func NewUser(name, email string) *User {
	return &User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
	}
}
