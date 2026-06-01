package entity

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint           `gorm:"primaryKey"                          json:"id"`
	FirstName string         `gorm:"size:200"                             json:"first_name"`
	LastName  string         `gorm:"size:200"                             json:"last_name"`
	Email     string         `gorm:"size:350"                             json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                json:"-"`
	Orders    []CustOrder    `gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE" json:"orders,omitempty"`
	Addresses []Address      `gorm:"many2many:customer_address;constraint:OnDelete:CASCADE" json:"addresses,omitempty"`
}

func (Customer) TableName() string {
	return "customer"
}