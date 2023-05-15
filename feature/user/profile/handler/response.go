package handler

import "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile"

type UserResponse struct {
	ID          string `json:"user_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

func CoreToUserResponse(data profile.Core) UserResponse {
	return UserResponse{
		ID:          data.ID,
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
	}
}
