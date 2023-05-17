package user

import (
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
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
	Position    admin.Position
	Office      admin.Office
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Submission struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"size:50"`
	TypeID    int    `gorm:"foreignKey:TypeID"`
	Title     string `gorm:"size:50;not null"`
	Message   string `gorm:"type:text;not null"`
	Status    string `gorm:"size:50;not null"`
	Is_Opened bool
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt
	Type      admin.Type
	Files     []File
	Tos       []To
	Ccs       []Cc
	Signs     []Sign
}

type File struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	SubmissionID int
	Name         string `gorm:"size:50;not null"`
	Link         string `gorm:"size:50;not null"`
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
}

type Sign struct {
	ID           int `gorm:"primaryKey"`
	SubmissionID int
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}
