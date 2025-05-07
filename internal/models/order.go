package models

import "time"

type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"constraint:OnDelete:CASCADE" json:"-"`
	Product   string    `gorm:"type:varchar(255);not null" json:"product"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (Order) TableName() string {
	return "orders"
}
