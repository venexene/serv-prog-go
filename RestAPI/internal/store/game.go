package store

import (
	"gravity-game-store/internal/entity"
	"gorm.io/gorm"
)

type GameStore struct{ db *gorm.DB }

func NewGameStore(d *gorm.DB) *GameStore { return &GameStore{db: d} }

func (s *GameStore) List(page, limit int) ([]entity.Game, int64, error) {
	var games []entity.Game
	var total int64
	if err := s.db.Model(&entity.Game{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	off := (page - 1) * limit
	if err := s.db.Preload("Authors").Offset(off).Limit(limit).Order("id ASC").Find(&games).Error; err != nil {
		return nil, 0, err
	}
	return games, total, nil
}

func (s *GameStore) ByID(id uint) (*entity.Game, error) {
	var g entity.Game
	if err := s.db.Preload("Authors").First(&g, id).Error; err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *GameStore) Create(g *entity.Game) error { return s.db.Create(g).Error }
func (s *GameStore) Update(g *entity.Game) error { return s.db.Save(g).Error }
func (s *GameStore) Delete(id uint) error         { return s.db.Delete(&entity.Game{}, id).Error }

func (s *GameStore) SetAuthors(g *entity.Game, ids []uint) error {
	var authors []entity.Author
	if len(ids) > 0 {
		if err := s.db.Find(&authors, ids).Error; err != nil {
			return err
		}
	}
	return s.db.Model(g).Association("Authors").Replace(authors)
}