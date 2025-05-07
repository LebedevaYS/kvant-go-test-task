package models

type User struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	Name     string  `gorm:"type:varchar(255);not null" json:"name"`
	Email    string  `gorm:"type:varchar(255);unique;not null" json:"email"`
	Age      int     `gorm:"not null" json:"age"`
	Password string  `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Orders   []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
}

func (User) TableName() string {
	return "users"
}
