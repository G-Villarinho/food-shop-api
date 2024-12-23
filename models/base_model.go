package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"column:Id;type:char(36);primaryKey"`
	CreatedAt time.Time      `gorm:"column:CreatedAt;not null"`
	UpdatedAt sql.NullTime   `gorm:"column:UpdatedAt;default:null"`
	DeletedAt gorm.DeletedAt `gorm:"column:DeletedAt;index"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	b.CreatedAt = time.Now().UTC()
	return
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = sql.NullTime{Time: time.Now().UTC(), Valid: true}
	return
}
