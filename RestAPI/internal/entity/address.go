package entity

import (
	"time"

	"gorm.io/gorm"
)

type Address struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	StreetNumber string         `gorm:"size:20"    json:"street_number"`
	StreetName   string         `gorm:"size:200"   json:"street_name"`
	City         string         `gorm:"size:200"   json:"city"`
	CountryID    int            `json:"country_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"       json:"-"`
}

func (Address) TableName() string {
	return "address"
}

type CustomerAddress struct {
	CustomerID uint `gorm:"primaryKey" json:"customer_id"`
	AddressID  uint `gorm:"primaryKey" json:"address_id"`
	StatusID   int  `json:"status_id"`
}

func (CustomerAddress) TableName() string {
	return "customer_address"
}