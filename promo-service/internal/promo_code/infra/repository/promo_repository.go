package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/iamrosada/microservice-goland/promo-service/internal/promo_code/entity"
	"gorm.io/gorm"
)

type PromoRepositoryImpl struct {
	DB *gorm.DB
}

// ApplyAll implements entity.PromotionRepository.
func (*PromoRepositoryImpl) ApplyAll(id uint) error {
	panic("unimplemented")
}

// ApplyToUsers implements entity.PromotionRepository.
func (*PromoRepositoryImpl) ApplyToUsers(id uint, userIds []int) error {
	panic("unimplemented")
}

// GetAppliedUsers implements entity.PromotionRepository.
func (*PromoRepositoryImpl) GetAppliedUsers(id uint) ([]uint, error) {
	panic("unimplemented")
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
		newID := generateNewID()

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

func generateNewID() uint {
	return uint(time.Now().UnixNano())
}
