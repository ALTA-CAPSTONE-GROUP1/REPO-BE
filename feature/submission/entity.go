package submission

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	FindRequirementHandler() echo.HandlerFunc
	AddSubmissionHandler() echo.HandlerFunc
	GetAllSubmissionHandler() echo.HandlerFunc
	GetSubmissionByIdHandler() echo.HandlerFunc
}

type UseCase interface {
	FindRequirementLogic(userID string, typeName string, value int) (Core, error)
	AddSubmissionLogic(newSub AddSubmissionCore, subFile *multipart.FileHeader) error
	GetAllSubmissionLogic(userID string, pr GetAllQueryParams) ([]AllSubmiisionCore, []SubTypeChoices, error)
	GetSubmissionByIDLogic(submissionID int, userID string) (GetSubmissionByIDCore, error)
}

type Repository interface {
	FindRequirement(userID string, typeName string, value int) (Core, error)
	InsertSubmission(newSub AddSubmissionCore) error
	SelectAllSubmissions(userID string, pr GetAllQueryParams) ([]AllSubmiisionCore, []SubTypeChoices, error)
	SelectSubmissionByID(submissionID int, userID string) (GetSubmissionByIDCore, error)
}

type AllSubmiisionCore struct {
	ID             int
	Tos            []ToApprover
	CCs            []CcApprover
	Title          string
	Status         string
	ReceiveDate    string
	Opened         bool
	Attachment     string
	SubmissionType string
}

type GetAllQueryParams struct {
	Title  string
	To     string
	Limit  int
	Offset int
}

type GetSubmissionByIDCore struct {
	To              []ToApprover
	Title           string
	CC              []CcApprover
	SubmissionType  string
	ApproverActions []ApproverActions
	ActionMessage   string
	Attachment      string
	Message         string
	Status          string
	Signs           []string
}

type ApproverActions struct {
	Action           string
	ApproverName     string
	ApproverPosition string
	Message          string
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

type SubTypeChoices struct {
	SubTypeName  string
	SubtypeValue []int
}
