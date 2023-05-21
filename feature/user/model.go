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
	User      admin.Users
	Files     []File `gorm:"foreignKey:SubmissionID"`
	Tos       []To   `gorm:"foreignKey:SubmissionID"`
	Ccs       []Cc   `gorm:"foreignKey:SubmissionID"`
	Signs     []Sign `gorm:"foreignKey:SubmissionID"`
}

type File struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	SubmissionID int
	Name         string     `gorm:"size:50;not null"`
	Link         string     `gorm:"size:50;not null"`
	Submission   Submission `gorm:"foreignKey:SubmissionID"`
}

type Cc struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	SubmissionID int
	UserID       string
	Name         string `gorm:"size:50;not null"`
	Is_Opened    bool
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	User         Users      `gorm:"foreignKey:UserID"`
	Submission   Submission `gorm:"foreignKey:SubmissionID"`
}

type To struct {
	ID           int `gorm:"primaryKey;autoIncrement"`
	SubmissionID int
	UserID       string
	Name         string `gorm:"size:50;not null"`
	Action_Type  string `gorm:"size:50;not null"`
	Is_Opened    bool
	Message      string
	CreatedAt    time.Time  `gorm:"autoCreateTime"`
	User         Users      `gorm:"foreignKey:UserID"`
	Submission   Submission `gorm:"foreignKey:SubmissionID"`
}

type Sign struct {
	ID           int       `gorm:"primaryKey"`
	SubmissionID int       `gorm:"foreignKey:SubmissionID"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	Name         string
	UserID       string
	User         Users      `gorm:"foreignKey:UserID"`
	Submission   Submission `gorm:"foreignKey:SubmissionID"`
}
