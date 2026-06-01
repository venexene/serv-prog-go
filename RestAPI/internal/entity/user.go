package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey"                          json:"id"`
	Username  string         `gorm:"uniqueIndex;size:100;not null"       json:"username"`
	Password  string         `gorm:"size:255;not null"                   json:"-"`
	Role      string         `gorm:"size:50;default:user"                json:"role"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                               json:"-"`
}

func (User) TableName() string {
	return "user"
}