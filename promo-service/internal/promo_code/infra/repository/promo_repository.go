package repository

import (
	"fmt"
	"strings"

	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/entity"
	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/infra/util"
	"gorm.io/gorm"
)

type PromoRepositoryImpl struct {
	DB *gorm.DB
}

func NewPromoRepository(db *gorm.DB) *PromoRepositoryImpl {
	return &PromoRepositoryImpl{
		DB: db,
	}
}

func (r *PromoRepositoryImpl) Create(promo *entity.Promotion) error {
	return r.DB.Create(promo).Error
}

func (r *PromoRepositoryImpl) FindByID(id uint) (*entity.Promotion, error) {
	var promo entity.Promotion
	err := r.DB.Where("id = ?", id).First(&promo).Error
	if err != nil {
		return nil, err
	}
	return &promo, nil
}

func (r *PromoRepositoryImpl) FindCodesByPromotionID(promotionID uint) ([]string, error) {
	var codes []string
	err := r.DB.Table("codes_promos").Select("Codes").Where("promotion_id = ?", promotionID).Pluck("Codes", &codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *PromoRepositoryImpl) FindAllPromos() ([]*entity.Promotion, error) {
	var promotions []*entity.Promotion
	err := r.DB.Table("promotions").Find(&promotions).Error
	if err != nil {
		return nil, err
	}
	return promotions, nil
}

func (r *PromoRepositoryImpl) AddCodes(id uint, codes []string) error {
	promo, err := r.FindByID(id)
	if err != nil {
		return err
	}

	codeMap := make(map[string]bool)
	for _, codeValue := range codes {
		lowercaseCode := strings.ToLower(codeValue)

		if _, exists := codeMap[lowercaseCode]; exists {
			return fmt.Errorf("duplicate code found: %s", codeValue)
		}

		codeMap[lowercaseCode] = true
		newID := util.GenerateNewID()

		code := entity.CodesPromo{
			ID:          newID,
			Codes:       []string{codeValue},
			PromotionID: promo.ID,
		}

		if err := r.DB.Create(&code).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *PromoRepositoryImpl) CreateUserThatWillUsePromo(promo *entity.UserPromotion) error {
	return r.DB.Create(promo).Error
}
