package repository

import (
	"errors"
	"fmt"

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
func (ar *approverModel) SelectSubmissionByHyperApproval(userID string, usrId string, id int, token string) (approve.GetSubmissionByIDCore, error) {
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
// UpdateUser implements approve.Repository
func (ar *approverModel) UpdateByHyperApproval(userID string, input approve.Core) error {
	var dbsub user.Submission
	var owner admin.Users
	var tos []user.To

	tx := ar.db.Model(&dbsub).
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN users ON tos.user_id = users.id").
		Where("tos.user_id = ? AND submissions.id = ?", input.UserID, input.ID).
		Find(&dbsub)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}
	fmt.Println(dbsub)

	tx = ar.db.Where("id = ?", dbsub.UserID).First(&owner)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}

	tx = ar.db.Where("submission_id = ?", dbsub.ID).Find(&tos)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}
	fmt.Println(tos)

	switch input.Status {
	case "approve":
		dbsub.Status = "Waiting"
	case "revise":
		dbsub.Status = "Revised"
	case "reject":
		dbsub.Status = "Rejected"
	default:
		return errors.New("invalid status")
	}

	if len(tos) > 0 && tos[len(tos)-1].UserID == input.UserID {
		dbsub.Status = "Approved"
	}

	tx = ar.db.Model(&dbsub).Updates(user.Submission{Status: dbsub.Status})

	if tx.RowsAffected == 0 {
		log.Error("no rows affected on update submission")
		return errors.New("data is up to date")
	}

	if tx.Error != nil {
		log.Error("error on update submission")
		return tx.Error
	}

	to := user.To{}
	tx = ar.db.Model(&user.To{}).
		Preload("User").
		Where("user_id = ? AND submission_id = ?", input.UserID, dbsub.ID).
		First(&to)

	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}

	tx = ar.db.Model(&to).
		Updates(user.To{Message: input.Message})

	// if tx.RowsAffected == 0 {
	// 	log.Error("no rows affected on update to message")
	// 	return errors.New("data is up to date")
	// }

	if tx.Error != nil {
		log.Error("error on update to message")
		return tx.Error
	}

	actionType := ""
	switch dbsub.Status {
	case "Approved":
		actionType = "approve"
	case "Waiting":
		actionType = "approve"
	case "Revised":
		actionType = "revise"
	case "Rejected":
		actionType = "reject"
	}

	tx = ar.db.Model(&user.To{}).
		Joins("JOIN users ON tos.user_id = users.id").
		Where("tos.user_id = ? AND tos.submission_id = ?", input.UserID, dbsub.ID).
		Update("action_type", actionType)

	if tx.Error != nil {
		log.Error("error on update action_type in 'to' table")
		return tx.Error
	}

	var sign string
	var signdb user.Sign
	if input.Status == "approve" {
		var err error
		sign, err = helper.GenerateUniqueSign(userID)
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
	}

	file := user.File{}
	if err := ar.db.Model(&user.File{}).Select("link").Where("submission_id = ?", dbsub.ID).First(&file).Error; err != nil {
		log.Error("no file found for the given submission ID")
		return errors.New("no file found")
	}

	recipient := []string{owner.Email}
	receiverName := []string{"Admin"}
	helper.SendSimpleEmail(input.Status, signdb.Name, dbsub.Title, "Update on your submission", recipient, receiverName, owner.Email)

	return nil
}
