package models

type Product struct {
	ID       int
	Name     string  `gorm:"not null"`
	Price    float64 `gorm:"type:numeric(10,2);not null;default:0"`
	Quantity int     `gorm:"not null;default:0"`
}
