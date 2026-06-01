package entity

import (
	"time"

	"gorm.io/gorm"
)

type Game struct {
	ID          uint           `gorm:"primaryKey"                          json:"id"`
	Title       string         `gorm:"size:400;not null"                   json:"title"`
	Genre       string         `gorm:"size:100"                             json:"genre"`
	Platform    string         `gorm:"size:100"                             json:"platform"`
	PublisherID int            `json:"publisher_id"`
	ReleaseDate string         `gorm:"size:10"                             json:"release_date"`
	NumPlayers  int            `json:"num_players"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                               json:"-"`
	Authors     []Author       `gorm:"many2many:game_author;constraint:OnDelete:CASCADE" json:"authors,omitempty"`
	OrderLines  []OrderLine    `gorm:"foreignKey:GameID;constraint:OnDelete:CASCADE"    json:"-"`
}

func (Game) TableName() string {
	return "game"
}

type GameAuthor struct {
	GameID   uint `gorm:"primaryKey"`
	AuthorID uint `gorm:"primaryKey"`
}

func (GameAuthor) TableName() string {
	return "game_author"
}