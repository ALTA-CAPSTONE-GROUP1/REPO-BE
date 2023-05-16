package submission

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	FindRequirementHandler() echo.HandlerFunc
	AddSubmissionHandler() echo.HandlerFunc
}

type UseCase interface {
	FindRequirementLogic(userID string, typeName string, value int) (Core, error)
	AddSubmissionLogic(newSub AddSubmissionCore, subFile *multipart.FileHeader) error
}

type Repository interface {
	FindRequirement(userID string, typeName string, value int) (Core, error)
	InsertSubmission(newSub AddSubmissionCore) error
}

type AddSubmissionCore struct {
	OwnerID           string
	ToApprover        []ToApprover
	CC                []CcApprover
	SubmissionType    string
	SubmissiontTypeID int
	Status            string
	SubmissionValue   int
	Title             string
	Message           string
	Attachment        string
	AttachmentLink    string
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
