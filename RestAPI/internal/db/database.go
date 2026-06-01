package db

import (
	"fmt"
	"os"
	"path/filepath"

	"gravity-game-store/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(path string, lvl logrus.Level) (*gorm.DB, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("mkdir %s: %w", dir, err)
	}

	gormLvl := logger.Silent
	if lvl <= logrus.DebugLevel {
		gormLvl = logger.Info
	}

	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(gormLvl),
	})
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("pragma: %w", err)
	}
	return db, nil
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&entity.User{},
		&entity.Author{},
		&entity.Game{},
		&entity.GameAuthor{},
		&entity.Customer{},
		&entity.Address{},
		&entity.CustomerAddress{},
		&entity.CustOrder{},
		&entity.OrderLine{},
		&entity.OrderHistory{},
		&entity.OrderStatus{},
		&entity.ShippingMethod{},
	)
}