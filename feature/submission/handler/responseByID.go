package handler

type ResponseByID struct {
	To             []ApproverRecipient `json:"to"`
	CC             []CCRecipient       `json:"cc"`
	SubmissionType string              `json:"submission_type"`
	Title          string              `json:"title"`
	Message        string              `json:"message"`
	ApproverAction []ApproverAction    `json:"approver_action"`
	ActionMessage  string              `json:"action_message"`
	Attachment     string              `json:"attachment"`
}

type ApproverRecipient struct {
	ApproverPosition string `json:"approver_position"`
	ApproverName     string `json:"approver_name"`
}

type CCRecipient struct {
	CCPosition string `json:"cc_position"`
	CCName     string `json:"cc_name"`
}

type ApproverAction struct {
	Action           string `json:"action"`
	ApproverName     string `json:"approver_name"`
	ApproverPosition string `json:"approver_position"`
}
