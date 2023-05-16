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
