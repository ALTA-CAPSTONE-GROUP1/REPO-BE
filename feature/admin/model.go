package admin

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	ID          string `gorm:"primaryKey;size:50"`
	OfficeID    int    `gorm:"foreignKey:OfficeID"`
	PositionID  int    `gorm:"foreignKey:PositionID"`
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
	ID       uint
	Name     string
	Level    string
	ParentID uint
	Parent   *Office `gorm:"foreignkey:ParentID"`
}

type Position struct {
	ID        int            `gorm:"primaryKey"`
	Name      string         `gorm:"size:50;not null"`
	Tag       string         `gorm:"size:50;not null"`
	Types     []Type         `gorm:"many2many:position_has_type;constraint:OnDelete:CASCADE;"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Type struct {
	ID          int            `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"size:50;not null"`
	Requirement string         `gorm:"size:255;not null"`
	Positions   []Position     `gorm:"many2many:position_has_type;constraint:OnDelete:CASCADE;"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type PositionHasType struct {
	ID         int    `gorm:"primaryKey;autoIncrement"`
	PositionID int    `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TypeID     int    `gorm:"primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	As         string `gorm:"size:10;not null"`
	ToLevel    int
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Value      int
	Position   Position `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Type       Type     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
