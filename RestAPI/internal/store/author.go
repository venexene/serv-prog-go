package store

import (
	"gravity-game-store/internal/entity"
	"gorm.io/gorm"
)

type AuthorStore struct{ db *gorm.DB }

func NewAuthorStore(d *gorm.DB) *AuthorStore { return &AuthorStore{db: d} }

func (s *AuthorStore) List(page, limit int) ([]entity.Author, int64, error) {
	var authors []entity.Author
	var total int64
	if err := s.db.Model(&entity.Author{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := (page - 1) * limit
	if err := s.db.Offset(off).Limit(limit).Order("id ASC").Find(&authors).Error; err != nil {
		return nil, 0, err
	}
	return authors, total, nil
}

func (s *AuthorStore) ByID(id uint) (*entity.Author, error) {
	var a entity.Author
	if err := s.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AuthorStore) WithGames(id uint) (*entity.Author, error) {
	var a entity.Author
	if err := s.db.Preload("Games").First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AuthorStore) Create(a *entity.Author) error { return s.db.Create(a).Error }
func (s *AuthorStore) Update(a *entity.Author) error { return s.db.Save(a).Error }
func (s *AuthorStore) Delete(id uint) error           { return s.db.Delete(&entity.Author{}, id).Error }