package submission

import "github.com/labstack/echo/v4"

type Handler interface {
	FindRequirementHandler() echo.HandlerFunc
}

type UseCase interface {
	FindRequirementLogic(userID string, typeName string, value int) (Core, error)
}

type Repository interface {
	FindRequirement(userID string, typeName string, value int) (Core, error)
}

type Core struct {
	To          []ToApprover
	CC          []CcApprover
	Requirement string
}

type ToApprover struct {
	ApproverPosition string
	ApproverId       string
	ApproverName     string
}

type CcApprover struct {
	CcPosition string
	CcName     string
	CcId       string
}
