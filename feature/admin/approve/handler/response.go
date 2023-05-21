package handler

type ResponseByID struct {
	Title          string           `json:"submission_title"`
	ApproverAction []ApproverAction `json:"approver_action"`
	Message        string           `json:"message_body"`
	SubmissionType string           `json:"submission_type"`
	Attachment     string           `json:"attachment"`
}

type ApproverAction struct {
	Action           string `json:"action"`
	ApproverName     string `json:"approver_name"`
	ApproverPosition string `json:"approver_position"`
}
