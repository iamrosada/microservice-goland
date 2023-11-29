package entity

type PromotionRepository interface {
	Create(promo *Promotion) error
	FindByID(id uint) (*Promotion, error)
	AddCodes(id uint, codes []string) error
	FindCodesByPromotionID(uint) ([]string, error)
	FindAllPromos() ([]*Promotion, error)
}
