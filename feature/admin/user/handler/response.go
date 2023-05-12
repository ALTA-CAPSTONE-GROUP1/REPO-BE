package handler

import "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"

type UserResponse struct {
	ID          string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	OfficeID    int    `json:"office_id"`
	PositionID  int    `json:"position_id"`
}

func CoreToUserResponse(data user.Core) UserResponse {
	return UserResponse{
		ID:          data.ID,
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
		OfficeID:    data.OfficeID,
		PositionID:  data.PositionID,
	}
}

func CoreToGetAllUserResponse(data []user.Core) []UserResponse {
	res := []UserResponse{}
	for _, val := range data {
		res = append(res, CoreToUserResponse(val))
	}
	return res
}
