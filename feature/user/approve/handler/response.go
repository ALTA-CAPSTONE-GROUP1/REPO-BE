package handler

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
)

type SubmissionResponse struct {
	ID        int         `json:"submission_id"`
	UserID    string      `json:"user_id"`
	TypeID    int         `json:"submission_type_name"`
	Title     string      `json:"title"`
	Message   string      `json:"message"`
	Status    string      `json:"status"`
	Is_Opened bool        `json:"is_opened"`
	Files     []user.File `json:"attachment"`
}

func CoreToApproveResponse(data submission.Core) SubmissionResponse {
	return SubmissionResponse{
		ID:          data.ID,
		Name:        data.Name,
		Email:       data.Email,
		PhoneNumber: data.PhoneNumber,
		Password:    data.Password,
		Position:    data.Position.Name,
		Office:      data.Office.Name,
	}
}

func CoreToGetAllApproveResponse(data []user.Core) []SubmissionResponse {
	res := make([]SubmissionResponse, len(data))
	for i, val := range data {
		res[i] = CoreToApproveResponse(val)
	}
	return res
}
