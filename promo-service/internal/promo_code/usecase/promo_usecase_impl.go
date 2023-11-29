package usecase

import (
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/entity"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/infra/util"
)

type PromotionUsecaseImpl struct {
	PromotionRepository entity.PromotionRepository
}

func NewPromoUsecase(repo entity.PromotionRepository) *PromotionUsecaseImpl {
	return &PromotionUsecaseImpl{
		PromotionRepository: repo,
	}
}

func (u *PromotionUsecaseImpl) CreatePromo(promotion *entity.Promotion) error {
	newID := util.GenerateNewID()
	promotion.ID = newID
	return u.PromotionRepository.Create(promotion)
}

func (u *PromotionUsecaseImpl) GetPromotionByID(id uint) (*entity.Promotion, error) {
	return u.PromotionRepository.FindByID(id)
}
func (u *PromotionUsecaseImpl) GetCodeByID(id uint) ([]string, error) {
	return u.PromotionRepository.FindCodesByPromotionID(id)
}

func (u *PromotionUsecaseImpl) AddCodesToPromotion(id uint, codes []string) error {

	return u.PromotionRepository.AddCodes(id, codes)
}

func (u *PromotionUsecaseImpl) GetAllPromos() ([]*entity.Promotion, error) {
	return u.PromotionRepository.FindAllPromos()
}
