package store

import (
	"gravity-game-store/internal/entity"
	"gorm.io/gorm"
)

type UserStore struct{ db *gorm.DB }

func NewUserStore(d *gorm.DB) *UserStore { return &UserStore{db: d} }

func (s *UserStore) ByUsername(name string) (*entity.User, error) {
	var u entity.User
	if err := s.db.Where("username = ?", name).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserStore) Create(u *entity.User) error { return s.db.Create(u).Error }