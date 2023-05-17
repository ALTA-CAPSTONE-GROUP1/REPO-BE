package handler

import (
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
)

type SubmissionResponse struct {
	ID             int    `json:"submission_id"`
	From           string `json:"user_id"`
	SubmissionType string `json:"submission_type"`
	Title          string `json:"title"`
	Message        string `json:"message"`
	Status         string `json:"status"`
	Is_Opened      bool   `json:"is_opened"`
	Attachment     string `json:"attachment"`
}

func CoreToApproveResponse(data approve.Core) SubmissionResponse {
	var fileLink []string
	for _, file := range data.Files {
		fileLink = append(fileLink, file.Link)
	}

	attachment := strings.Join(fileLink, ", ")

	return SubmissionResponse{
		ID:             data.ID,
		From:           data.UserID,
		SubmissionType: data.Type.Name,
		Title:          data.Title,
		Message:        data.Message,
		Status:         data.Status,
		Is_Opened:      false,
		Attachment:     attachment,
	}
}

func CoreToGetAllApproveResponse(data []approve.Core) []SubmissionResponse {
	res := make([]SubmissionResponse, len(data))
	for i, val := range data {
		res[i] = CoreToApproveResponse(val)
	}
	return res
}
