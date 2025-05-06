package models

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null"`
	User      User      `gorm:"constraint:OnDelete:CASCADE"`
	Product   string    `gorm:"type:varchar(255);not null"`
	Quantity  int       `gorm:"not null"`
	Price     float64   `gorm:"type:decimal(10,2);not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

func (Order) TableName() string {
	return "orders"
}
