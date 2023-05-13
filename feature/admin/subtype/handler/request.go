package handler

type AddSubTypeRequest struct {
	SubmissionTypeName string     `json:"submission_type_name" validate:"required"`
	Position           []string   `json:"position" validate:"required"`
	SubmissionValues   []subValue `json:"submission_value" validate:"required"`
	Requirement        string     `json:"requirement" validate:"required"`
}

type subValue struct {
	Value      int      `json:"value" validate:"required"`
	PositionTo []string `json:"position_to" validate:"required"`
	PositionCC []string `json:"position_cc" validate:"required"`
}
