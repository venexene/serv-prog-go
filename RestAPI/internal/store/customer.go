package store

import (
	"gravity-game-store/internal/entity"
	"gorm.io/gorm"
)

type CustomerStore struct{ db *gorm.DB }

func NewCustomerStore(d *gorm.DB) *CustomerStore { return &CustomerStore{db: d} }

func (s *CustomerStore) List(page, limit int) ([]entity.Customer, int64, error) {
	var cust []entity.Customer
	var total int64
	if err := s.db.Model(&entity.Customer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := (page - 1) * limit
	if err := s.db.Offset(off).Limit(limit).Order("id ASC").Find(&cust).Error; err != nil {
		return nil, 0, err
	}
	return cust, total, nil
}

func (s *CustomerStore) ByID(id uint) (*entity.Customer, error) {
	var c entity.Customer
	if err := s.db.First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CustomerStore) WithOrders(id uint) (*entity.Customer, error) {
	var c entity.Customer
	if err := s.db.Preload("Orders.OrderLines.Game").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *CustomerStore) Create(c *entity.Customer) error { return s.db.Create(c).Error }
func (s *CustomerStore) Update(c *entity.Customer) error { return s.db.Save(c).Error }
func (s *CustomerStore) Delete(id uint) error             { return s.db.Delete(&entity.Customer{}, id).Error }