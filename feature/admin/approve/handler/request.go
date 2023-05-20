package handler

type InputUpdate struct {
	SubID  int    `json:"submission_id" form:"submission_id"`
	Action string `json:"new_status" form:"new_status"`
}

type InputGet struct {
	SubID int    `json:"submission_id" form:"submission_id"`
	Token string `json:"token" form:"token"`
}
