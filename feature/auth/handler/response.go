package handler

type LoginResponse struct {
	Role  string `json:"role"`
	Token string `json:"token"`
}

type SignResponse struct {
	SubmissionTitle  string `json:"submission_title"`
	OfficialName     string `json:"official_name"`
	OfficialPosition string `json:"official_position"`
	Date             string `json:"date"`
}
