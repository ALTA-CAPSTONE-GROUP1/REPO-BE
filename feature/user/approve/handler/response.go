package handler

import (
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
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

type SubmissionByIdResponse struct {
	ID             int           `json:"submission_id"`
	From           string        `json:"from"`
	Title          string        `json:"title"`
	SubmissionType string        `json:"submission_type"`
	Message        string        `json:"message"`
	Attachment     string        `json:"attachment"`
	To             []ToApp       `json:"to"`
	Cc             []CcRecipient `json:"cc"`
	Status         []Action      `json:"status_by"`
}

type ToApp struct {
	ToName     string `json:"name"`
	ToPosition string `json:"position"`
}

type CcRecipient struct {
	CcName     string `json:"name"`
	CcPosition string `json:"position"`
}

type Action struct {
	Action    string `json:"status"`
	AppAction string `json:"by"`
}

func CoreToApproveByIdResponse(data approve.Core) SubmissionByIdResponse {
	result := SubmissionByIdResponse{
		ID:             data.ID,
		From:           data.UserID,
		Title:          data.Title,
		Message:        data.Message,
		SubmissionType: data.Type.Name,
	}
	for _, v := range data.Tos {
		cTos := ToApp{
			ToName:     v.User.Name,
			ToPosition: v.User.Position.Name,
		}
		result.To = append(result.To, cTos)
	}

	for _, y := range data.Ccs {
		cCcs := CcRecipient{
			CcName:     y.User.Name,
			CcPosition: y.User.Position.Name,
		}
		result.Cc = append(result.Cc, cCcs)
	}

	return result
}

func SubmissionToCore(data user.Submission) approve.Core {
	result := approve.Core{
		ID:        data.ID,
		UserID:    data.UserID,
		TypeID:    data.TypeID,
		Title:     data.Title,
		Message:   data.Message,
		Status:    data.Status,
		Is_Opened: false,
		CreatedAt: time.Time{},
		Type:      admin.Type{Name: data.Type.Name},
		User:      admin.Users{Name: data.User.Name, Position: data.User.Position},
	}

	for _, v := range data.Tos {
		cTos := user.To{
			User: user.Users{
				Position: v.User.Position,
				Name:     v.User.Name,
			},
		}
		result.Tos = append(result.Tos, cTos)
	}

	for _, y := range data.Ccs {
		cCcs := user.Cc{
			User: user.Users{
				Position: y.User.Position,
				Name:     y.User.Name,
			},
		}
		result.Ccs = append(result.Ccs, cCcs)
	}

	return result
}
