package handler

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
)

type UserResponse struct {
	ID          string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Position    string `json:"position"`
	Office      string `json:"office"`
}

func CoreToUserResponse(data user.Core) UserResponse {
	return UserResponse{
		ID:          data.ID,
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
		Position:    data.Position.Name,
		Office:      data.Office.Name,
	}
}

func CoreToGetAllUserResponse(data []user.Core) []UserResponse {
	res := make([]UserResponse, len(data))
	for i, val := range data {
		res[i] = CoreToUserResponse(val)
	}
	return res
}
