package entity

type UserPromotion struct {
	ID          uint `json:"id"`
	UserID      uint `json:"user_id"`
	PromotionID uint `json:"promotion_id"`
	Type        int  `json:"type"`
}

type PromotionRepository interface {
	GetAvailableUsers(promoType int) ([]int, error)
	GetAppliedUsers(promoID uint) ([]int, error)
	ApplyPromotion(promoID uint, userIDs []int) error
}
