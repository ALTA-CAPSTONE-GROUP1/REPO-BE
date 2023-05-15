package handler

type InputUpdate struct {
	Email       string `json:"email" form:"email"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Password    string `json:"password" form:"password"`
}
