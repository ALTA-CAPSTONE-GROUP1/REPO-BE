package handler

type RegisterInput struct {
	OfficeID    int    `json:"office_id"`
	PositionID  int    `json:"position_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}
