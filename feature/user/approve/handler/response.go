package handler

import (
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
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
		To: make([]ToApp, len(data.Tos)),
		Cc: make([]CcRecipient, len(data.Ccs)),
	}

	fromUser := data.User
	result.From = FromApp{
		Name:     fromUser.Name,
		Position: fromUser.Position.Name,
	}

	result.Title = data.Title
	result.Message = data.Message
	result.SubmissionType = data.Type.SubmissionTypeName
	// result.Attachment = data.Attachment

	for i, v := range data.Tos {
		toUser := v.User
		result.To[i] = ToApp{
			ToName:     toUser.Name,
			ToPosition: toUser.Position.Name,
		}
	}

	for i, y := range data.Ccs {
		ccUser := y.User
		result.Cc[i] = CcRecipient{
			CcName:     ccUser.Name,
			CcPosition: ccUser.Position.Name,
		}
	}

	// for _, z := range data.StatusBy {
	// 	result.StatusBy = append(result.StatusBy, Action{
	// 		Action:    z.Status,
	// 		AppAction: z.By,
	// 	})
	// }

	return result
}

func SubmissionToCore(data aMod.Submission) approve.Core {
	result := approve.Core{
		ID:        data.ID,
		UserID:    data.UserID,
		TypeID:    data.TypeID,
		Title:     data.Title,
		Message:   data.Message,
		Status:    data.Status,
		Is_Opened: false,
		CreatedAt: data.CreatedAt,
		Type:      subtype.Core{SubmissionTypeName: data.Type.Name},
	}

	result.User.Name = data.User.Name
	result.User.Position.Name = data.User.Position.Name

	for _, v := range data.Tos {
		cTos := approve.ToCore{
			User: user.Core{
				Name: v.User.Name,
				Position: position.Core{
					Name: v.User.Position.Name,
				},
			},
		}
		result.Tos = append(result.Tos, cTos)
	}

	for _, y := range data.Ccs {
		cCcs := approve.CcCore{
			User: user.Core{
				Name: y.User.Name,
				Position: position.Core{
					Name: y.User.Position.Name,
				},
			},
		}
		result.Ccs = append(result.Ccs, cCcs)
	}

	for _, z := range data.Signs {
		cSigns := approve.SignCore{
			User: user.Core{
				Name: z.User.Name,
				Position: position.Core{
					Name: z.User.Position.Name,
				},
			},
		}
		result.Signs = append(result.Signs, cSigns)
	}

	return result
}
