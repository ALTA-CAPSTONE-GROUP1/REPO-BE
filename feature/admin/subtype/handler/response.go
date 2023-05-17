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

type Meta struct {
	CurrentLimit  int `json:"current_limit"`
	CurrentOffset int `json:"current_offset"`
	CurrentPage   int `json:"current_page"`
	TotalData     int `json:"total_data"`
	TotalPage     int `json:"total_page"`
}
