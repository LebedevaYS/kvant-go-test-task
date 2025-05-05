package models

import "time"

type Order struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    UserID    uint      `json:"user_id"`
    Product   string    `json:"product"`
    Quantity  int       `json:"quantity"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
