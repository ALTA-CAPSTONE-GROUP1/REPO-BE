package submission

type Handler interface {
}

type UseCase interface {
}

type Repository interface {
	FindRequirement(userID string, typeName string, value string) (RequirementDB, error)
}

type Approver struct {
	ApproverPositionID int
	ApproverPosition   string
	ApproverDetail     []ApproverDetail
	Level              int
}

type CC struct {
	CCPositionID string
	CCPosition   string
	CCDetails    string
}

type RequirementDB struct {
	ApplicantID         string
	ApplicantName       string
	ApplicantPositionID string
	ApplicantPosition   string
	ApplicantOfficeID   string
	ApplicantOfficeName int
	To                  []Approver
	CC                  []CC
	Requirement         string
}

type ApproverDetail struct {
	ApproverID         string
	ApproverName       string
	ApproverOfficeID   string
	ApproverOfficeName string
}

type CcDetail struct {
	CCID         string
	CCName       string
	CCOfficeID   string
	CCOfficeName string
}
