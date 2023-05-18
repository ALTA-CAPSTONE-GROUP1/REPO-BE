package repository

import (
	"time"

	"gorm.io/gorm"
)

type Submission struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"size:50"`
	TypeID    int    `gorm:"not null"`
	Title     string `gorm:"size:50;not null;unique"`
	Message   string `gorm:"type:text;not null"`
	Status    string `gorm:"size:50;not null"`
	Is_Opened bool
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
	Files     []File
	Tos       []To
	Ccs       []Cc
	Signs     []Sign
}

type File struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	SubmissionID int
	Name         string `gorm:"size:50;not null"`
	Link         string `gorm:"size:255;not null"`
}

type Cc struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	SubmissionID int
	UserID       string `gorm:"size:50"`
	Name         string `gorm:"size:50;not null"`
	Is_Opened    bool
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

type To struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	SubmissionID int
	UserID       string `gorm:"size:50"`
	Name         string `gorm:"size:50;not null"`
	Action_Type  string `gorm:"size:50;not null"`
	Is_Opened    bool
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Message      string
}

type Sign struct {
	ID           int `gorm:"primaryKey"`
	SubmissionID int
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Name         string    `gorm:"size:50"`
	UserID       string    `gorm:"size:50"`
}
