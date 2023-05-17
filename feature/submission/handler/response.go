package handler

type RequirementResponseBody struct {
	To          []ToApprover `json:"to"`
	CC          []CcApprover `json:"cc"`
	Requirement string       `json:"requirement"`
}

type CcApprover struct {
	CcPosition string `json:"cc_position"`
	CcName     string `json:"cc_name"`
	CcId       string `json:"cc_id"`
}

type ToApprover struct {
	ApproverPosition string `json:"approver_position"`
	ApproverId       string `json:"approver_id"`
	ApproverName     string `json:"approver_name"`
}
type SubmissionResponse struct {
	Submissions           []Submission           `json:"submissions"`
	SubmissionTypeChoices []SubmissionTypeChoice `json:"submission_type_choices"`
}
type Submission struct {
	ID             int        `json:"id"`
	To             []Approver `json:"to"`
	CC             []CC       `json:"cc"`
	Title          string     `json:"title"`
	Status         string     `json:"status"`
	Attachment     string     `json:"attachment"`
	ReceiveDate    string     `json:"receive_date"`
	Opened         bool       `json:"opened"`
	SubmissionType string     `json:"submission_type"`
}

type Approver struct {
	ApproverPosition string `json:"approver_position"`
	ApproverName     string `json:"approver_name"`
}

type CC struct {
	CCPosition string `json:"cc_position"`
	CCName     string `json:"cc_name"`
}

type SubmissionTypeChoice struct {
	Name   string `json:"name"`
	Values []int  `json:"values"`
}

type Meta struct {
	CurrentLimit  int `json:"current_limit"`
	CurrentOffset int `json:"current_offset"`
	CurrentPage   int `json:"current_page"`
	TotalData     int `json:"total_data"`
	TotalPage     int `json:"total_page"`
}