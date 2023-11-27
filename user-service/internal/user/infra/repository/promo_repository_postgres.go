package repository

import (
	"fmt"
	"math/rand"

	"github.com/iamrosada/microservice-goland/user-service/internal/user/entity"
	"gorm.io/gorm"
)

type PromotionRepositoryPostgres struct {
	DB *gorm.DB
}

func NewPromotionRepositoryPostgres(db *gorm.DB) *PromotionRepositoryPostgres {
	return &PromotionRepositoryPostgres{DB: db}
}

func (r *PromotionRepositoryPostgres) GetAvailableUsers(promoType int) ([]int, error) {
	// 1. Consultar o banco de dados para obter IDs de usuários que receberam a promoção
	var usersPromotion []entity.UserPromotion
	if err := r.DB.Where("type != ?", promoType).Find(&usersPromotion).Error; err != nil {
		return nil, err
	}

	// 2. Consultar o banco de dados para obter todos os IDs de usuários
	var allUsers []entity.User
	if err := r.DB.Find(&allUsers).Error; err != nil {
		return nil, err
	}

	// 3. Identificar os IDs dos usuários que não receberam a promoção
	var usersWithoutPromotionIDs []int
	for _, user := range allUsers {
		if !containsUserPromotion(usersPromotion, int(user.ID)) {
			usersWithoutPromotionIDs = append(usersWithoutPromotionIDs, int(user.ID))
		}
	}

	// 4. Verificar se há usuários disponíveis
	if len(usersWithoutPromotionIDs) == 0 {
		// Não há usuários disponíveis
		return nil, nil
	}

	// 5. Retornar a lista final de usuários que ainda podem receber a promoção
	return usersWithoutPromotionIDs, nil
}

// containsUserPromotion verifica se um usuário já recebeu uma determinada promoção
func containsUserPromotion(usersPromotion []entity.UserPromotion, userID int) bool {
	for _, userPromotion := range usersPromotion {
		if int(userPromotion.UserID) == userID {
			return true
		}
	}
	return false
}
func (r *PromotionRepositoryPostgres) GetAppliedUsers(promoID uint) ([]int, error) {
	// Consultar o banco de dados para obter IDs de usuários aos quais o código promocional foi aplicado
	var appliedUsers []entity.UserPromotion
	if err := r.DB.Find(&appliedUsers, "promotion_id = ?", promoID).Error; err != nil {
		return nil, err
	}

	// Verificar se nenhum usuário foi encontrado
	if len(appliedUsers) == 0 {
		return nil, fmt.Errorf("no users found with promotion applied for promoID %d", promoID)
	}

	// Extrair os IDs dos usuários
	var appliedUserIDs []int
	for _, appliedUser := range appliedUsers {
		appliedUserIDs = append(appliedUserIDs, int(appliedUser.UserID))
	}

	return appliedUserIDs, nil
}

func (r *PromotionRepositoryPostgres) ApplyPromotion(promoID uint, userIDs []int) error {
	// Iterate over user IDs and apply the promotion to each user
	randomNumber := rand.Intn(9) + 1

	for _, userID := range userIDs {
		appliedPromotion := entity.UserPromotion{
			PromotionID: promoID,
			UserID:      uint(userID),
			Type:        int(randomNumber),
		}

		// Save the applied promotion to the database
		if err := r.DB.Create(&appliedPromotion).Error; err != nil {
			return err
		}
	}

	return nil
}
