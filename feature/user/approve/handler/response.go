package handler

import (
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	aMod "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	// aRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/users/approve/repository"
)

type SubmissionResponse struct {
	ID             int    `json:"submission_id"`
	Title          string `json:"title"`
	From           string `json:"from"`
	SubmissionType string `json:"submission_type"`
	Status         string `json:"status"`
	CreatedAt      string `json:"receive_date"`
	Is_Opened      bool   `json:"opened"`
}

func CoreToApproveResponse(data approve.Core) SubmissionResponse {
	return SubmissionResponse{
		ID:             data.ID,
		Title:          data.Title,
		From:           data.UserID,
		SubmissionType: data.Type.SubmissionTypeName,
		Status:         data.Status,
		CreatedAt:      data.CreatedAt.Add(7 * time.Hour).Format("2006-01-02 15:04"),
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
		Name:     data.Owner.Name,
		Position: data.Owner.Position,
	}

	result.Title = data.Title
	result.Message = data.Message
	result.SubmissionType = data.Type.SubmissionTypeName

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

	for _, v := range data.Tos {
		tmp := Action{
			Action:    v.Action_Type,
			AppAction: v.User.ID,
		}
		result.StatusBy = append(result.StatusBy, tmp)
	}

	result.Attachment = data.Files[0].Link

	return result
}

func SubmissionToCore(usr admin.Users, file []aMod.File, receivers []admin.Users, ccs []admin.Users, subData aMod.Submission) approve.Core {
	var res approve.Core
	var tos []approve.ToCore
	var ccsCore []approve.CcCore
	var sign []approve.SignCore
	var owner approve.OwnerCore
	var coreFile []approve.FileCore

	for _, v := range receivers {
		tos = append(tos, approve.ToCore{
			SubmissionID: subData.ID,
			Name:         v.Name,
			Position:     v.Position.Name,
			UserID:       v.ID,
		})
	}

	for _, v := range ccs {
		ccsCore = append(ccsCore, approve.CcCore{
			SubmissionID: subData.ID,
			UserID:       v.ID,
			Name:         v.Name,
			Position:     v.Position.Name,
		})
	}

	for _, v := range file {
		coreFile = append(coreFile, approve.FileCore{
			SubmissionID: subData.ID,
			Name:         v.Name,
			Link:         v.Link,
		})
	}

	owner = approve.OwnerCore{
		SubmissionID: subData.ID,
		Name:         usr.Name,
		Position:     usr.Position.Name,
	}

	res = approve.Core{
		ID:        subData.ID,
		UserID:    subData.UserID,
		TypeID:    subData.TypeID,
		Title:     subData.Title,
		Message:   subData.Message,
		Status:    subData.Status,
		Is_Opened: false,
		CreatedAt: subData.CreatedAt,
		Type:      subtype.Core{SubmissionTypeName: subData.Type.Name},
		Files:     coreFile,
		Tos:       tos,
		Ccs:       ccsCore,
		Signs:     sign,
		Owner:     owner,
	}

	return res
}
