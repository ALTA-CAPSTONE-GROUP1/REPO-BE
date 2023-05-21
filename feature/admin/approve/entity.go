package approve

import (
	"time"

	cType "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	cUser "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
	"github.com/labstack/echo/v4"
)

type Core struct {
	ID        int
	UserID    string
	TypeID    int
	Title     string
	Message   string
	Status    string
	Is_Opened bool
	CreatedAt time.Time
	Type      cType.Core
	User      cUser.Core
	Files     []FileCore
	Tos       []ToCore
	Ccs       []CcCore
	Signs     []SignCore
}

type FileCore struct {
	ID           int
	SubmissionID int
	Name         string
	Link         string
	Submission   Core
}

type CcCore struct {
	ID           int
	SubmissionID int
	UserID       string
	Name         string
	Is_Opened    bool
	CreatedAt    time.Time
	// User         cUser.Core
	Submission Core
	Position   string
	User       string
}

type ToCore struct {
	ID           int
	SubmissionID int
	UserID       string
	Name         string
	Action_Type  string
	Is_Opened    bool
	Message      string
	CreatedAt    time.Time
	// User         cUser.Core
	Submission Core
	Position   string
	User       string
}

type SignCore struct {
	ID           int
	SubmissionID int
	CreatedAt    time.Time
	Name         string
	UserID       string
}

type Handler interface {
	GetSubmissionByHyperApprovalHandler() echo.HandlerFunc
	UpdateByHyperApprovalHandler() echo.HandlerFunc
}

type UseCase interface {
	GetSubmissionByHyperApproval(userID string, id int, token string) (Core, error)
	UpdateByHyperApproval(userID string, updateInput Core) error
}

type Repository interface {
	SelectSubmissionByHyperApproval(userID string, id int, token string) (Core, error)
	UpdateByHyperApproval(userID string, input Core) error
}
