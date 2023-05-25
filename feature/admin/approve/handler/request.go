package handler

type InputUpdate struct {
	SubID  int    `json:"submission_id" form:"submission_id"`
	UsrID  string `json:"user_id" form:"user_id"`
	Action string `json:"new_status" form:"new_status"`
}

type InputGet struct {
	SubID int    `json:"submission_id"`
	UsrID string `json:"user_id" form:"user_id"`
	Token string `json:"token"`
}
