package repository

import (
	"fmt"
	"log"

	"github.com/iamrosada/microservice-goland/user-service/internal/user/entity"
	"github.com/iamrosada/microservice-goland/user-service/internal/user/infra/util"
	"gorm.io/gorm"
)

type PromotionRepositoryPostgres struct {
	DB *gorm.DB
}

func NewPromotionRepositoryPostgres(db *gorm.DB) *PromotionRepositoryPostgres {
	return &PromotionRepositoryPostgres{DB: db}
}

func (r *PromotionRepositoryPostgres) GetAvailableUsers(promoType int) ([]int, error) {
	var usersPromotion []entity.UserPromotion
	if err := r.DB.Where("type != ?", promoType).Find(&usersPromotion).Error; err != nil {
		return nil, err
	}

	var allUsers []entity.User
	if err := r.DB.Find(&allUsers).Error; err != nil {
		return nil, err
	}

	var usersWithoutPromotionIDs []int
	for _, user := range allUsers {
		if !util.ContainsUserPromotion(usersPromotion, int(user.ID)) {
			usersWithoutPromotionIDs = append(usersWithoutPromotionIDs, int(user.ID))
		}
	}

	if len(usersWithoutPromotionIDs) == 0 {
		return nil, nil
	}

	return usersWithoutPromotionIDs, nil
}

func (r *PromotionRepositoryPostgres) GetAppliedUsers(promoID uint) ([]int, error) {
	var appliedUsers []entity.UserPromotion
	if err := r.DB.Find(&appliedUsers, "promotion_id = ?", promoID).Error; err != nil {
		return nil, err
	}

	if len(appliedUsers) == 0 {
		return nil, fmt.Errorf("no users found with promotion applied for promoID %d", promoID)
	}

	var appliedUserIDs []int
	for _, appliedUser := range appliedUsers {
		appliedUserIDs = append(appliedUserIDs, int(appliedUser.UserID))
	}

	return appliedUserIDs, nil
}

func (r *PromotionRepositoryPostgres) ApplyPromotion(promoID uint, userIDs []int) error {
	promoTypeFromOtherMicroservice, err := util.FetchPromoTypeFromMicroservice(fmt.Sprintf("http://localhost:8080/promo/%d", promoID))
	if err != nil {
		log.Printf("Error fetching promo type: %v", err)
		return fmt.Errorf("failed to fetch promo type: %v", err)
	}

	log.Printf("Promo type from microservice: %d", promoTypeFromOtherMicroservice)

	for _, userID := range userIDs {
		appliedPromotion := entity.UserPromotion{
			ID:          util.GenerateNewID(),
			PromotionID: promoID,
			UserID:      uint(userID),
			Type:        promoTypeFromOtherMicroservice,
		}

		log.Printf("Applying promotion for user %d with promo type %d", userID, promoTypeFromOtherMicroservice)

		if err := r.DB.Create(&appliedPromotion).Error; err != nil {
			log.Printf("Error applying promotion for user %d: %v", userID, err)
			return fmt.Errorf("failed to apply promotion for user %d: %v", userID, err)
		}
	}

	return nil
}
