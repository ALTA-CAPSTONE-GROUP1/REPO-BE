package approve

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
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
	Files     []user.File
	Tos       []user.To
	Ccs       []user.Cc
	Signs     []user.Sign
}

// type Core struct {
// 	To          []ToApprover
// 	CC          []CcApprover
// 	Requirement string
// }

// type ToApprover struct {
// 	ApproverPosition string
// 	ApproverId       string
// 	ApproverName     string
// }

// type CcApprover struct {
// 	CcPosition string
// 	CcName     string
// 	CcId       string
// }

type Handler interface {
	GetSubmissionAprroveHandler() echo.HandlerFunc
}

type UseCase interface {
	GetSubmissionAprrove(limit, offset int, search string) ([]Core, error)
}

type Repository interface {
	SelectSubmissionAprrove(limit, offset int, search string) ([]Core, error)
}
