package repository

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          string         `gorm:"primaryKey;size:50"`
	OfficeID    int            `gorm:"foreignKey:OfficeID"`
	PositionID  int            `gorm:"foreignKey:PositionID"`
	Name        string         `gorm:"size:50;not null"`
	Email       string         `gorm:"size:50"`
	PhoneNumber string         `gorm:"size:50"`
	Password    string         `gorm:"size:50"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
