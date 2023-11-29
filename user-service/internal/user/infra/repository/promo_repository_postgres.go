package repository

import (
	"fmt"
	"math/rand"
	"strconv"

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
		if !containsUserPromotion(usersPromotion, int(user.ID)) {
			usersWithoutPromotionIDs = append(usersWithoutPromotionIDs, int(user.ID))
		}
	}

	if len(usersWithoutPromotionIDs) == 0 {
		return nil, nil
	}

	return usersWithoutPromotionIDs, nil
}

func containsUserPromotion(usersPromotion []entity.UserPromotion, userID int) bool {
	for _, userPromotion := range usersPromotion {
		if int(userPromotion.UserID) == userID {
			return true
		}
	}
	return false
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
	randomNumber := rand.Intn(9) + 1

	// users, err := fetchUsersFromUserAppliedMicroservice(fmt.Sprintf("http://localhost:8000/api/users/promo/%s/applied", id))
	// log.Printf("users: %+v", users)

	// if err != nil {
	// 	log.Printf("Error fetching users from user microservice: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users from user microservice"})
	// 	return
	// }

	// var userIDs []int
	// for _, user := range users {
	// 	userIDs = append(userIDs, int(user))
	// }

	for _, userID := range userIDs {
		appliedPromotion := entity.UserPromotion{
			ID:          util.GenerateNewID(),
			PromotionID: promoID,
			UserID:      uint(userID),
			Type:        int(randomNumber),
		}
		fmt.Printf("%s ", strconv.FormatUint(uint64(userID), 10))

		if err := r.DB.Create(&appliedPromotion).Error; err != nil {
			return err
		}
	}

	return nil
}

func fetchUsersFromUserAppliedMicroservice(s string) {
	panic("unimplemented")
}
