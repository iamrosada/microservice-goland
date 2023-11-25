package entity

type Promotion struct {
	ID          uint
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Type        int    `json:"type"`
	// Codes       []Code `gorm:"foreignKey:PromotionID"`
}
