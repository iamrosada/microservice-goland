package entity

type Promotion struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Type        int    `json:"type"`
}
