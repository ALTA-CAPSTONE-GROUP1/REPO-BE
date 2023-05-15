package handler

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
)

type UserResponse struct {
	ID          string         `json:"user_id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	Password    string         `json:"password"`
	Position    admin.Position `json:"position"`
	Office      admin.Office   `json:"office"`
}

func CoreToUserResponse(data user.Core) UserResponse {
	return UserResponse{
		ID:          data.ID,
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
		Position: admin.Position{
			Name: data.Position.Name,
		},
		Office: admin.Office{
			Name: data.Office.Name,
		},
	}
}

func CoreToGetAllUserResponse(data []user.Core) []UserResponse {
	res := []UserResponse{}
	for _, val := range data {
		res = append(res, CoreToUserResponse(val))
	}
	return res
}
