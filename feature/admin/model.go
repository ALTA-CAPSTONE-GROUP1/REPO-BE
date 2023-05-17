package admin

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID          string `gorm:"primaryKey;size:50"`
	OfficeID    int
	PositionID  int
	Name        string `gorm:"size:50;not null"`
	Email       string `gorm:"size:50"`
	PhoneNumber string `gorm:"size:50"`
	Password    string `gorm:"size:100"`
	Position    Position
	Office      Office
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Office struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:50;not null"`
}

type Position struct {
	ID        int            `gorm:"primaryKey"`
	Name      string         `gorm:"size:50;not null"`
	Tag       string         `gorm:"size:50;not null"`
	Types     []Type         `gorm:"many2many:position_has_types;constraint:OnDelete:CASCADE;"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Type struct {
	ID          int            `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"size:50;not null;unique"`
	Requirement string         `gorm:"size:255;not null"`
	Positions   []Position     `gorm:"many2many:position_has_types;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type PositionHasType struct {
	ID         int `gorm:"primaryKey;autoIncrement"`
	PositionID int
	TypeID     int
	As         string `gorm:"size:10;not null"`
	ToLevel    int
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Value      int
	Position   Position `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type       Type     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
