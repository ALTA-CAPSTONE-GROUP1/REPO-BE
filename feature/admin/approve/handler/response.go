package handler

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
)

type SubmissionByHyperApprovalResponse struct {
	Title          string         `json:"submission-title"`
	AppName        string         `json:"applicant_name"`
	AppPosition    string         `json:"applicant_position"`
	To             []ToApp        `json:"approver_action"`
	Message        string         `json:"message_body"`
	SubmissionType string         `json:"submission_type"`
	Attachment     []AttachmentAp `json:"attachment"`
}

type ToApp struct {
	ToAction   string `json:"action"`
	ToName     string `json:"approver_name"`
	ToPosition string `json:"approver_position"`
}

type AttachmentAp struct {
	// Name string
	Link string `json:"attachment"`
}

func CoreToSubmissionByHyperApprovalResponse(data approve.Core) SubmissionByHyperApprovalResponse {
	result := SubmissionByHyperApprovalResponse{
		Title:          data.Title,
		AppName:        data.User.Name,
		AppPosition:    data.User.Position.Name,
		Message:        data.Message,
		SubmissionType: data.Type.SubmissionTypeName,
	}

	for _, v := range data.Tos {
		tmp := ToApp{
			ToAction:   v.Action_Type,
			ToName:     v.Name,
			ToPosition: v.Position,
		}
		result.To = append(result.To, tmp)
	}

	for _, v := range data.Files {
		tmp := AttachmentAp{
			// Name: v.Name,
			Link: v.Link,
		}
		result.Attachment = append(result.Attachment, tmp)
	}

	return result
}

// ===========================================================================
// ===========================================================================
// ===========================================================================

// func SubmissionToCore(submision user.Submission) approve.Core {
// 	result := approve.Core{
// 		ID:        submision.ID,
// 		UserID:    submision.UserID,
// 		TypeID:    submision.TypeID,
// 		Title:     submision.Title,
// 		Message:   submision.Message,
// 		Status:    submision.Message,
// 		Is_Opened: false,
// 		CreatedAt: submision.CreatedAt,
// 		// Type:      subtype.Core{},
// 		// User:      user.Core{},
// 		// Files:     []approve.FileCore{},
// 		// Tos:       []approve.ToCore{},
// 		// Ccs:       []approve.CcCore{},
// 		// Signs:     []approve.SignCore{},
// 	}

// 	if !reflect.ValueOf(submision.User).IsZero() {
// 		result.User = cUser.Core{
// 			Name:     submision.User.Name,
// 			Position: cPos.Core{},
// 		}
// 	}

// 	if !reflect.ValueOf(product.Category).IsZero() {
// 		result.Category = categories.CategoryEntity{
// 			Category: product.Category.Category,
// 		}
// 	}

// 	for _, v := range submission.ToApprover {
// 		var image = productImages.ProductImagesEntity{
// 			Id:    v.ID,
// 			Image: v.Image,
// 		}
// 		result.ProductImage = append(result.ProductImage, image)
// 	}

// 	return result
// }

// func ListProductToProductEntity(product []Product) []products.ProductEntity {
// 	var productEntity []products.ProductEntity
// 	for _, v := range product {
// 		productEntity = append(productEntity, ProductToProductEntity(v))
// 	}
// 	return productEntity
// }
