package handler

import (
	"fmt"
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	aMod "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	// aRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/users/approve/repository"
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
		SubmissionType: data.Type.SubmissionTypeName,
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

// ===========================================================================================
// ===========================================================================================
// ===========================================================================================
// ===========================================================================================

type SubmissionByIdResponse struct {
	ID             int           `json:"submission_id"`
	From           FromApp       `json:"from"`
	To             []ToApp       `json:"to"`
	Cc             []CcRecipient `json:"cc"`
	Title          string        `json:"title"`
	SubmissionType string        `json:"submission_type"`
	StatusBy       []Action      `json:"status_by"`
	Message        string        `json:"message"`
	Attachment     string        `json:"attachment"`
}

type FromApp struct {
	Name     string `json:"name"`
	Position string `json:"position"`
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
		ID: data.ID,
	}

	result.From = FromApp{
		Name:     data.User.Name,
		Position: data.User.Position.Name,
	}

	result.Title = data.Title
	result.Message = data.Message
	result.SubmissionType = data.Type.SubmissionTypeName
	// result.Attachment = data.Attachment

	for _, v := range data.Tos {
		tmp := ToApp{
			ToName:     v.Name,
			ToPosition: v.Position,
		}
		result.To = append(result.To, tmp)
	}

	for _, v := range data.Ccs {
		tmp := CcRecipient{
			CcName:     v.Name,
			CcPosition: v.Position,
		}
		result.Cc = append(result.Cc, tmp)
	}

	// for _, z := range data.StatusBy {
	// 	result.StatusBy = append(result.StatusBy, Action{
	// 		Action:    z.Status,
	// 		AppAction: z.By,
	// 	})
	// }

	return result
}

func SubmissionToCore(rcv []admin.Users, ccs []admin.Users, submissiondatas aMod.Submission) approve.Core {
	var res approve.Core
	var tos []approve.ToCore
	var ccsCore []approve.CcCore

	for _, v := range rcv {
		tos = append(tos, approve.ToCore{
			SubmissionID: submissiondatas.ID,
			Name:         v.Name,
			Position:     v.Position.Name,
			UserID:       v.ID,
		})
	}

	for _, v := range ccs {
		ccsCore = append(ccsCore, approve.CcCore{
			SubmissionID: submissiondatas.ID,
			UserID:       v.ID,
			Name:         v.Name,
			Position:     v.Position.Name,
		})
	}
	res = approve.Core{
		ID:        submissiondatas.ID,
		UserID:    submissiondatas.UserID,
		TypeID:    submissiondatas.TypeID,
		Title:     submissiondatas.Title,
		Message:   submissiondatas.Message,
		Status:    submissiondatas.Status,
		Is_Opened: false,
		CreatedAt: submissiondatas.CreatedAt,
		Tos:       tos,
		Ccs:       ccsCore,
	}
	fmt.Println(tos)
	fmt.Println(ccsCore)
	return res
}
