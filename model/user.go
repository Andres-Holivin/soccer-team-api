package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`

	CreatedAt time.Time      `gorm:"autoCreateTime;type:timestamptz" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;type:timestamptz" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
