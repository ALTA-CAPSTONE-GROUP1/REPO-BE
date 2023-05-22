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
	Owner     OwnerCore
	StatusBy  []StatusBy
	Approv    ApproverCore
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
	User         cUser.Core
	Submission   Core
	Position     string
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
	User         cUser.Core
	Submission   Core
	Position     string
}

type SignCore struct {
	ID           int
	SubmissionID int
	CreatedAt    time.Time
	Name         string
	UserID       string
	User         cUser.Core
	Submission   Core
}

type OwnerCore struct {
	SubmissionID int
	Name         string
	Position     string
}

type ApproverCore struct {
	SubmissionID int
	Action       string
}

type StatusBy struct {
	Action   string
	Position string
}

type GetAllQueryParams struct {
	FromTo string
	Title  string
	Type   string
	Limit  int
	Offset int
}

type Handler interface {
	GetSubmissionAprroveHandler() echo.HandlerFunc
	GetSubmissionByIdHandler() echo.HandlerFunc
	UpdateSubmissionApproveHandler() echo.HandlerFunc
}

type UseCase interface {
	GetSubmissionAprrove(userID string, search GetAllQueryParams) ([]Core, error)
	GetSubmissionById(userID string, id int) (Core, error)
	UpdateApprove(userID string, id int, updateInput Core) error
}

type Repository interface {
	SelectSubmissionAprrove(userID string, search GetAllQueryParams) ([]Core, error)
	SelectSubmissionById(userID string, id int) (Core, error)
	UpdateApprove(userID string, id int, input Core) error
}
