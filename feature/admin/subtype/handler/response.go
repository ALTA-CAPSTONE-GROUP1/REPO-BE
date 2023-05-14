package handler

type GetSubmissionTypeResponse struct {
	SubmissionType []SubmissionType `json:"submission_type"`
	Position       []Position       `json:"positions"`
}

type Position struct {
	PositionName string `json:"position_name"`
	PositionTag  string `json:"position_tag"`
}

type SubmissionType struct {
	SubmissionTypeName string             `json:"submission_type_name"`
	SubmissionDetail   []SubmissionDetail `json:"submission_detail"`
}

type SubmissionDetail struct {
	SubmissionValue       int    `json:"submission_value"`
	SubmissionRequirement string `json:"submission_requirement"`
}
