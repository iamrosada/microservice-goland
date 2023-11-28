package entity

type Serializable interface {
	BeforeSave() error
	AfterFind() error
}

type Code struct {
	ID          uint     `json:"id"`
	Codes       []string `gorm:"type:jsonb" json:"codes"`
	PromotionID uint     `json:"promotion_id"`
}
