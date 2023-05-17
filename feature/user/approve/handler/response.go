package handler

import (
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
)

type SubmissionResponse struct {
	ID             int       `json:"submission_id"`
	Title          string    `json:"title"`
	From           string    `json:"from"`
	SubmissionType string    `json:"submission_type"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"receive_date"`
	Is_Opened      bool      `json:"opened"`
}

func CoreToApproveResponse(data approve.Core) SubmissionResponse {
	return SubmissionResponse{
		ID:             data.ID,
		Title:          data.Title,
		From:           data.UserID,
		SubmissionType: data.Type.Name,
		Status:         data.Status,
		CreatedAt:      time.Time{},
		Is_Opened:      false,
	}
}

func CoreToGetAllApproveResponse(data []approve.Core) []SubmissionResponse {
	res := make([]SubmissionResponse, len(data))
	for i, val := range data {
		res[i] = CoreToApproveResponse(val)
	}
	return res
}
