package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
	sRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/repository"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
)

type approverModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) approve.Repository {
	return &approverModel{
		db: db,
	}
}

// / SelectSubmissionByHyperApproval implements approve.Repository
func (ar *approverModel) SelectSubmissionByHyperApproval(userID string, id int, token string) (approve.GetSubmissionByIDCore, error) {
	var (
		result         approve.GetSubmissionByIDCore
		submissionByID sRepo.Submission
		subTypeDetails admin.Type
	)

	if token != config.TokenSuperAdmin {
		return approve.GetSubmissionByIDCore{}, errors.New("invalid token")
	}

	if err := ar.db.Where("id = ?", id).
		Preload("Files").
		Preload("Tos").
		Preload("Ccs").
		Preload("Signs").
		First(&submissionByID).Error; err != nil {
		log.Errorf("error on finding submissions for by submissionid %s:", err)
		return approve.GetSubmissionByIDCore{}, err
	}

	if err := ar.db.Where("id = ?", submissionByID.TypeID).Find(&subTypeDetails).Error; err != nil {
		log.Errorf("error in finding submissionType Name %w", err)
		return approve.GetSubmissionByIDCore{}, err
	}

	// var toApprover []submission.ToApprover
	var toActions []approve.ApproverActions
	for _, to := range submissionByID.Tos {
		var toDetails admin.Users
		if err := ar.db.Where("id = ?", to.UserID).Preload("Position").First(&toDetails).Error; err != nil {
			log.Errorf("failed on finding positions of tos %w", err)
			return approve.GetSubmissionByIDCore{}, err
		}
		toActions = append(toActions, approve.ApproverActions{
			Action:           to.Action_Type,
			ApproverName:     toDetails.Name,
			ApproverPosition: toDetails.Position.Name,
		})

	}

	if len(submissionByID.Files) > 0 {
		result.Attachment = submissionByID.Files[0].Link
	}

	result.ApproverActions = append(result.ApproverActions, toActions...)
	result.Title = submissionByID.Title
	result.Message = submissionByID.Message
	result.Status = submissionByID.Status
	result.SubmissionType = subTypeDetails.Name

	return result, nil
}

// UpdateUser implements approve.Repository
func (ar *approverModel) UpdateByHyperApproval(userID string, input approve.Core) error {
	var dbsub user.Submission
	var owner admin.Users
	var tos []user.To

	if err := ar.db.Model(&user.Submission{}).Where("id = ?", input.ID).First(&dbsub).Error; err != nil {
		log.Error("no rows found for the given submission ID")
		return errors.New("no data found")
	}

	if err := ar.db.Where("id = ?", dbsub.UserID).First(&owner).Error; err != nil {
		log.Error("no rows found for the given user and submission ID")
		return errors.New("no data found")
	}

	if err := ar.db.Where("submission_id = ?", dbsub.ID).Find(&tos).Error; err != nil {
		log.Error("no rows found for the given user and submission ID")
		return errors.New("no data found")
	}

	switch input.Status {
	case "approve":
		input.Status = "Approved"
	case "revise":
		input.Status = "Revised"
	case "reject":
		input.Status = "Rejected"
	default:
		return errors.New("invalid status")
	}

	// if len(tos) > 0 && tos[len(tos)-1] == userID {
	// 	dbsub.Status = "Approved"
	// }

	if err := ar.db.Model(&dbsub).Where("id = ?", input.ID).Updates(user.Submission{Status: input.Status}).Error; err != nil {
		log.Error("no rows affected on update submission")
		return errors.New("data is up to date")
	}

	to := user.To{}
	if err := ar.db.Model(&user.To{}).Where("submission_id = ?", input.ID).First(&to).Error; err != nil {
		log.Error("no rows found for the given submission ID in user.To")
		return errors.New("no data found in user.To")
	}

	// if to.ID != 0 {
	// 	if err := ar.db.Model(&user.To{}).Where("submission_id = ?", input.ID).Updates(approve.ToCore{Message: "Action from Admin, thank you!"}).Error; err != nil {
	// 		log.Error("error on update to message")
	// 		return err
	// 	}
	// }

	// actionType := ""
	// switch input.Status {
	// case "waiting":
	// 	actionType = "approve"
	// case "revised":
	// 	actionType = "revise"
	// case "rejected":
	// 	actionType = "reject"
	// }

	// if err := ar.db.Model(&user.To{}).
	// 	// Joins("JOIN users ON user_id = users.id").
	// 	Where("submission_id = ?", input.ID).
	// 	Update("action_type", actionType).Error; err != nil {
	// 	log.Error("error on update action_type in 'to' table")
	// 	return err
	// }

	sign, err := helper.GenerateUniqueSign(userID)
	if err != nil {
		log.Error("failed to generate unique sign")
		return err
	}

	signdb := user.Sign{
		UserID:       userID,
		Name:         sign,
		SubmissionID: dbsub.ID,
	}

	if err := ar.db.Create(&signdb).Error; err != nil {
		log.Error("error on update sign")
		return err
	}

	file := user.File{}
	if err := ar.db.Model(&user.File{}).Select("link").Where("submission_id = ?", dbsub.ID).First(&file).Error; err != nil {
		log.Error("no file found for the given submission ID")
		return errors.New("no file found")
	}

	recipient := []string{owner.Email}
	receiverName := []string{"Admin"}
	helper.SendSimpleEmail(signdb.Name, dbsub.Title, "Update on your submission", recipient, receiverName, owner.Email)

	return nil
}
