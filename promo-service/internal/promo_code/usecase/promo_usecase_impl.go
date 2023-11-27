package usecase

import "github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/entity"

type PromotionUsecaseImpl struct {
	PromotionRepository entity.PromotionRepository
}

func NewPromoUsecase(repo entity.PromotionRepository) *PromotionUsecaseImpl {
	return &PromotionUsecaseImpl{
		PromotionRepository: repo,
	}
}

func (u *PromotionUsecaseImpl) CreatePromo(Promotion *entity.Promotion) error {
	return u.PromotionRepository.Create(Promotion)
}

func (u *PromotionUsecaseImpl) GetPromotionByID(id uint) (*entity.Promotion, error) {
	return u.PromotionRepository.FindByID(id)
}

func (u *PromotionUsecaseImpl) AddCodesToPromotion(id uint, codes []string) error {
	return u.PromotionRepository.AddCodes(id, codes)
}

func (u *PromotionUsecaseImpl) ApplyPromotionToAll(id uint) error {
	return u.PromotionRepository.ApplyAll(id)
}

func (u *PromotionUsecaseImpl) ApplyPromotionToUsers(id uint, userIds []int) error {
	return u.PromotionRepository.ApplyToUsers(id, userIds)
}

func (u *PromotionUsecaseImpl) GetAppliedUsers(id uint) ([]uint, error) {
	return u.PromotionRepository.GetAppliedUsers(id)
}
