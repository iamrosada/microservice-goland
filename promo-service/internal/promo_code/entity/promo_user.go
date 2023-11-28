package entity

type UserPromotion struct {
	ID          uint `json:"id"`
	UserID      uint `json:"user_id"`
	PromotionID uint `json:"promotion_id"`
	Type        int  `json:"type"`
}
