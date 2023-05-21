package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
	sRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/repository"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
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
		return approve.GetSubmissionByIDCore{}, errors.New("Invalid token")
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

// // UpdateUser implements approve.Repository
// func (ar *approverModel) UpdateByHyperApproval(userID string, input approve.Core) error {
// 	// submission := approve.Core{}
// 	dbsub := user.Submission{}
// 	tx := ar.db.Model(&user.Submission{}).
// 		Where("id = ?", input.ID).
// 		First(&dbsub)

// 	if tx.RowsAffected == 0 {
// 		log.Error("no rows found for the given submission ID")
// 		return errors.New("no data found")
// 	}

// 	switch input.Status {
// 	case "approve":
// 		input.Status = "Waiting"
// 	case "revise":
// 		input.Status = "Revised"
// 	case "reject":
// 		input.Status = "Rejected"
// 	default:
// 		return errors.New("invalid status")
// 	}

// 	tx = ar.db.Model(&dbsub).
// 		Where("id = ?", input.ID).
// 		Updates(user.Submission{Status: input.Status})

// 	if tx.RowsAffected == 0 {
// 		log.Error("no rows affected on update submission")
// 		return errors.New("data is up to date")
// 	}

// 	if tx.Error != nil {
// 		log.Error("error on update submission")
// 		return tx.Error
// 	}

// 	to := user.To{}
// 	tx = ar.db.Model(&user.To{}).
// 		Where("submission_id = ?", input.ID).
// 		First(&to)

// 	if tx.RowsAffected == 0 {
// 		log.Error("no rows found for the given submission ID in user.To")
// 		return errors.New("no data found in user.To")
// 	}

// 	tx = ar.db.Model(&user.To{}).
// 		Where("submission_id = ?", input.ID).
// 		Updates(approve.ToCore{Message: config.AutoMessageHyperApp})

// 	if tx.RowsAffected == 0 {
// 		log.Error("no rows affected on update to message")
// 		return errors.New("data is up to date")
// 	}

// 	if tx.Error != nil {
// 		log.Error("error on update to message")
// 		return tx.Error
// 	}

// 	actionType := ""
// 	switch input.Status {
// 	case "waiting":
// 		actionType = "approve"
// 	case "revised":
// 		actionType = "revise"
// 	case "rejected":
// 		actionType = "reject"
// 	}

// 	tx = ar.db.Model(&user.To{}).
// 		Joins("JOIN users ON user_id = users.id").
// 		Where("submission_id = ?", input.ID).
// 		Update("action_type", actionType)

// 	if tx.Error != nil {
// 		log.Error("error on update action_type in 'to' table")
// 		return tx.Error
// 	}

//		return nil
//	}
func (ar *approverModel) UpdateByHyperApproval(userID string, input approve.Core) error {
	dbsub := user.Submission{}
	tx := ar.db.Model(&user.Submission{}).
		Where("id = ?", input.ID).
		First(&dbsub)

	if tx.RowsAffected == 0 {
		log.Error("no rows found for the given submission ID")
		return errors.New("no data found")
	}

	switch input.Status {
	case "approve":
		input.Status = "Waiting"
	case "revise":
		input.Status = "Revised"
	case "reject":
		input.Status = "Rejected"
	default:
		return errors.New("invalid status")
	}

	tx = ar.db.Model(&dbsub).
		Where("id = ?", input.ID).
		Updates(user.Submission{Status: input.Status})

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
		Where("submission_id = ?", input.ID).
		First(&to)

	if tx.RowsAffected == 0 {
		log.Error("no rows found for the given submission ID in user.To")
		return errors.New("no data found in user.To")
	}

	if to.Action_Type == "" {
		userID := to.UserID

		userWithPosition := user.Users{}
		tx = ar.db.Model(&user.Users{}).
			Where("id = ?", userID).
			Preload("Position").
			First(&userWithPosition)

		if tx.RowsAffected == 0 {
			log.Error("no rows found for the given user ID in user.Users")
			return errors.New("no data found in user.Users")
		}

		positionHasType := admin.PositionHasType{}
		tx = ar.db.Model(&admin.PositionHasType{}).
			Where("type_id = ? AND position_id = ?", input.TypeID, userWithPosition.Position.ID).
			First(&positionHasType)

		if tx.RowsAffected == 0 {
			log.Error("no rows found for the given type and position in admin.PositionHasType")
			return errors.New("no data found in admin.PositionHasType")
		}

		// cari user id dengan toLevel terkecil dari position_has_type
		smallestToLevelUserID := ""
		tx = ar.db.Raw("SELECT user_id FROM user.To WHERE action_type IS NOT NULL ORDER BY action_type LIMIT 1").
			Scan(&smallestToLevelUserID)

		if tx.Error != nil {
			log.Error("error on retrieving the user ID with the smallest position_has_type to_level")
			return tx.Error
		}

		// update action_type di user.to tabel yg toLevel terkecil dari user id
		tx = ar.db.Model(&user.To{}).
			Where("submission_id = ? AND user_id = ?", input.ID, smallestToLevelUserID).
			Update("action_type", positionHasType.ToLevel)

		if tx.RowsAffected == 0 {
			log.Error("no rows affected on update to message")
			return errors.New("data is up to date")
		}

		if tx.Error != nil {
			log.Error("error on update to message")
			return tx.Error
		}
	}

	return nil
}
