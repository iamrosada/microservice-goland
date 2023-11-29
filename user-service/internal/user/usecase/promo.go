package usecase

import "github.com/iamrosada/microservice-goland/user-service/internal/user/entity"

// PromotionUseCase represents the use case for promotion operations
type PromotionUseCase struct {
	repo entity.PromotionRepository
}

// NewPromotionUseCase creates a new instance of PromotionUseCase
func NewPromotionUseCase(repo entity.PromotionRepository) *PromotionUseCase {
	return &PromotionUseCase{repo: repo}
}

// GetAvailableUsers retrieves available users for a given promotion type
func (uc *PromotionUseCase) GetAvailableUsers(promoType int) ([]int, error) {
	return uc.repo.GetAvailableUsers(promoType)

}

// GetAppliedUsers retrieves users to which the promotion has been applied
func (uc *PromotionUseCase) GetAppliedUsers(promoID uint) ([]int, error) {
	return uc.repo.GetAppliedUsers(promoID)
}

// ApplyPromotion applies the promotion to the specified users
func (uc *PromotionUseCase) ApplyPromotion(promoID uint, userIDs []int) error {
	return uc.repo.ApplyPromotion(promoID, userIDs)

}
