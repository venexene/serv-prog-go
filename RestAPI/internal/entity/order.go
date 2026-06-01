package entity

import (
	"time"

	"gorm.io/gorm"
)

type CustOrder struct {
	ID               uint           `gorm:"primaryKey"                          json:"id"`
	OrderDate        time.Time      `json:"order_date"`
	CustomerID       uint           `json:"customer_id"`
	ShippingMethodID int            `json:"shipping_method_id"`
	DestAddressID    uint           `json:"dest_address_id"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index"                                json:"-"`
	Customer         Customer       `gorm:"foreignKey:CustomerID"                json:"customer,omitempty"`
	OrderLines       []OrderLine    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_lines,omitempty"`
	OrderHistories   []OrderHistory `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
}

func (CustOrder) TableName() string {
	return "cust_order"
}

type OrderLine struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	OrderID uint    `json:"order_id"`
	GameID  uint    `json:"game_id"`
	Price   float64 `gorm:"type:decimal(10,2)" json:"price"`
	Game    Game    `gorm:"foreignKey:GameID"   json:"game,omitempty"`
}

func (OrderLine) TableName() string {
	return "order_line"
}

type OrderHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `json:"order_id"`
	StatusID   int       `json:"status_id"`
	StatusDate time.Time `json:"status_date"`
}

func (OrderHistory) TableName() string {
	return "order_history"
}

type OrderStatus struct {
	StatusID    int    `gorm:"primaryKey" json:"status_id"`
	StatusValue string `gorm:"size:50"    json:"status_value"`
}

func (OrderStatus) TableName() string {
	return "order_status"
}

type ShippingMethod struct {
	MethodID   int     `gorm:"primaryKey" json:"method_id"`
	MethodName string  `gorm:"size:100"   json:"method_name"`
	Cost       float64 `gorm:"type:decimal(10,2)" json:"cost"`
}

func (ShippingMethod) TableName() string {
	return "shipping_method"
}