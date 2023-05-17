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
	ID           int    `gorm:"primaryKey;autoIncrement"`
	SubmissionID int    `gorm:"foreignKey:SubmissionID"`
	UserID       string `gorm:"foreignKey:UserID"`
	Name         string `gorm:"size:50;not null"`
	Is_Opened    bool
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	User         Users
}

type To struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	SubmissionID int    `gorm:"foreignKey:SubmissionID"`
	UserID       string `gorm:"foreignKey:UserID"`
	Name         string `gorm:"size:50;not null"`
	Action_Type  string `gorm:"size:50;not null"`
	Is_Opened    bool
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	User         Users
}

type Sign struct {
	ID           int `gorm:"primaryKey"`
	SubmissionID int
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// func SubmissionToCore(data Submission) approve.Core {
// 	result := approve.Core{
// 		ID:        data.ID,
// 		UserID:    data.UserID,
// 		TypeID:    data.TypeID,
// 		Title:     data.Title,
// 		Message:   data.Message,
// 		Status:    data.Status,
// 		Is_Opened: false,
// 		CreatedAt: time.Time{},
// 		Type:      admin.Type{Name: data.Type.Name},
// 		User:      admin.Users{Name: data.User.Name, Position: data.User.Position},
// 	}

// 	for _, v := range data.Tos {
// 		cTos := To{
// 			User: Users{
// 				Position: v.User.Position,
// 				Name:     v.User.Name,
// 			},
// 		}
// 		result.Tos = append(result.Tos, cTos)
// 	}

// 	for _, y := range data.Ccs {
// 		cCcs := Cc{
// 			User: Users{
// 				Position: y.User.Position,
// 				Name:     y.User.Name,
// 			},
// 		}
// 		result.Ccs = append(result.Ccs, cCcs)
// 	}

// 	return result
// }

// for _, v := range dbsub {
// 	tmp := approve.Core{
// 		ID:        v.ID,
// 		UserID:    v.UserID,
// 		TypeID:    v.TypeID,
// 		Title:     v.Title,
// 		Message:   v.Message,
// 		Status:    v.Status,
// 		Is_Opened: false,
// 		CreatedAt: time.Time{},
// 		Type:      admin.Type{Name: v.Type.Name},
// 		User:      admin.Users{Name: v.User.Name, Position: v.User.Position},
// 		Files:     []user.File{},
// 		Tos: []user.To{
// 			user.To{
// 				User: user.Users{
// 					Name:     v.User.Name,
// 					Position: v.User.Position,
// 				},
// 			},
// 		},
// 		Ccs: []user.Cc{
// 			user.Cc{
// 				User: user.Users{
// 					Name:     v.User.Name,
// 					Position: v.User.Position,
// 				},
// 			},
// 		},
// 		Signs: []user.Sign{},
// 	}
// 	res = append(res, tmp)
// }
