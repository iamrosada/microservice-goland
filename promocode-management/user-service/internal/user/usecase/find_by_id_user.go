package usecase

import "github.com/iamrosada/microservice-goland/user-service/internal/user/entity"

type GetUserByIDInputDto struct {
	ID string `json:"id"`
}

type GetUserByIDOutputDto struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUserByIDUseCase struct {
	UserRepository entity.UserRepository
}

func NewGetUserByIDUseCase(UserRepository entity.UserRepository) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{UserRepository: UserRepository}
}

func (u *GetUserByIDUseCase) Execute(input GetUserByIDInputDto) (*GetUserByIDOutputDto, error) {
	User, err := u.UserRepository.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserByIDOutputDto{
		ID:   User.ID,
		Name: User.Name,

		Email: User.Email,
	}, nil
}
