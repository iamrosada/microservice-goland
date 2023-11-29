package usecase

import "github.com/iamrosada/microservice-goland/user-service/internal/user/entity"

type PromotionUseCase struct {
	repo entity.PromotionRepository
}

func NewPromotionUseCase(repo entity.PromotionRepository) *PromotionUseCase {
	return &PromotionUseCase{repo: repo}
}

func (uc *PromotionUseCase) GetAvailableUsers(promoType int) ([]int, error) {
	return uc.repo.GetAvailableUsers(promoType)

}

func (uc *PromotionUseCase) GetAppliedUsers(promoID uint) ([]int, error) {
	return uc.repo.GetAppliedUsers(promoID)
}

func (uc *PromotionUseCase) ApplyPromotion(promoID uint, userIDs []int) error {
	return uc.repo.ApplyPromotion(promoID, userIDs)

}
