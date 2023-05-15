package handler

type RegisterInput struct {
	OfficeID    int    `json:"office_id"`
	PositionID  int    `json:"position_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type InputUpdate struct {
	Name        string `json:"name" form:"name"`
	Email       string `json:"email" form:"email"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
	Password    string `json:"password" form:"password"`
	OfficeID    int    `json:"office_id" form:"office_id"`
	PositionID  int    `json:"position_id"  form:"position_id"`
}
