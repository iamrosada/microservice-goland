// infra/promo_repository.go
package repository

import (
	"fmt"
	"math/rand"
	"strings"

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
	return r.DB.Create(&promo).Error
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

	// Executando uma query SQL personalizada para selecionar os valores da coluna "Codes"
	err := r.DB.Raw("SELECT Codes FROM codes_promos WHERE promotion_id = ?", promotionID).Pluck("Codes", &codes).Error
	if err != nil {
		return nil, err
	}

	return codes, nil
}

func (r *PromoRepositoryImpl) FindAllCodes() ([]entity.CodesPromo, error) {
	var codes []entity.CodesPromo
	err := r.DB.Find(&codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}

func (r *PromoRepositoryImpl) AddCodes(id uint, codes []string) error {
	// Find the promotion by ID
	promo, err := r.FindByID(id)
	if err != nil {
		return err
	}

	// Assuming that there is a one-to-many relationship between Promotion and CodesPromo
	// Create a map to check for duplicate codes
	codeMap := make(map[string]bool)

	// Check for duplicates (case-insensitive) and add each unique code to the promotion
	for _, codeValue := range codes {
		// Convert the code to lowercase for case-insensitive comparison
		lowercaseCode := strings.ToLower(codeValue)

		// Check if the code is already in the map (duplicate)
		if _, exists := codeMap[lowercaseCode]; exists {
			// Handle duplicate code error
			return fmt.Errorf("duplicate code found: %s", codeValue)
		}

		// Add the lowercase code to the map
		codeMap[lowercaseCode] = true
		newID := generateNewID()

		// Assign the generated ID to the promotion
		// Create the CodesPromo entity
		code := entity.CodesPromo{
			ID:          newID,
			Codes:       []string{codeValue},
			PromotionID: promo.ID, // Assuming Promotion has an ID field
		}

		// Save the code to the database
		if err := r.DB.Create(&code).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *PromoRepositoryImpl) CreateUserThatWillUsePromo(promo *entity.UserPromotion) error {
	return r.DB.Create(&promo).Error
}

func generateNewID() uint {

	return uint(rand.Uint32())
}
