package entity

type PromotionRepository interface {
	Create(promo *Promotion) error
	FindByID(id uint) (*Promotion, error)
	AddCodes(id uint, codes []string) error
	ApplyAll(id uint) error
	ApplyToUsers(id uint, userIds []int) error
	GetAppliedUsers(id uint) ([]uint, error)
	FindCodesByPromotionID(uint) ([]string, error)
}
