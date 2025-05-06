package models

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"type:varchar(255);unique;not null"`
	Age      int    `gorm:"not null"`
	Password string `gorm:"column:password_hash;type:varchar(255);not null"`
}

func (User) TableName() string {
	return "users"
}
