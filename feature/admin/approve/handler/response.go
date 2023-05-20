package handler

type SubmissionByHyperApprovalResponse struct {
	Title          string         `json:"submission-title"`
	AppName        string         `json:"applicant_name"`
	AppPosition    string         `json:"applicant_position"`
	To             []ToApp        `json:"approver_action"`
	Message        string         `json:"message_body"`
	SubmissionType string         `json:"submission_type"`
	Attachment     []AttachmentAp `json:"attachment"`
}

type ToApp struct {
	ToAction   string `json:"action"`
	ToName     string `json:"approver_name"`
	ToPosition string `json:"approver_position"`
}

type AttachmentAp struct {
	Link string `json:"attachment"`
}
