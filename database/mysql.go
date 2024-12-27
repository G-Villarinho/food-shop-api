package database

import (
	"context"

	"github.com/G-Villarinho/food-shop-api/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMysqlConnection(ctx context.Context) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.Env.ConnectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	slqDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := slqDB.PingContext(ctx); err != nil {
		_ = slqDB.Close()
		return nil, err
	}

	return db, nil
}
