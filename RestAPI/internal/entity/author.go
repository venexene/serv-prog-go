package entity

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID        uint           `gorm:"primaryKey"                          json:"id"`
	Name      string         `gorm:"size:400;not null"                   json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                               json:"-"`
	Games     []Game         `gorm:"many2many:game_author;constraint:OnDelete:CASCADE" json:"games,omitempty"`
}

func (Author) TableName() string {
	return "author"
}