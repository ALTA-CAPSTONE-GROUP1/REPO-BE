package handler

type LoginInput struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type SignVaidation struct{
	SignID string `json:"sign_id"`
}