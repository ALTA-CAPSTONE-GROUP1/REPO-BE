package handler

type InputUpdate struct {
	Action  string `json:"action" form:"action"`
	Message string `json:"approval_message" form:"approval_message"`
}
