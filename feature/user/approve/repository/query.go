package repository

import (
	"errors"
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	uMod "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve/handler"
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

// UpdateUser implements approve.Repository
func (ar *approverModel) UpdateApprove(userID string, id int, input approve.Core) error {
	var submission user.Submission

	tx := ar.db.Model(&submission).
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN users ON tos.user_id = users.id").
		Where("tos.user_id = ? AND submissions.id = ?", userID, id).
		Find(&submission)

	if tx.RowsAffected == 0 {
		log.Error("no rows found for the given user and submission ID")
		return errors.New("no data found")
	}

	switch input.Status {
	case "approve":
		submission.Status = "waiting"
	case "revise":
		submission.Status = "revised"
	case "reject":
		submission.Status = "rejected"
	default:
		return errors.New("invalid status")
	}

	tx = ar.db.Model(&submission).Updates(user.Submission{Status: submission.Status})

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
		Where("user_id = ? AND submission_id = ?", userID, submission.ID).
		First(&to)

	if tx.RowsAffected == 0 {
		log.Error("no rows found for the given user and submission ID in user.To")
		return errors.New("no data found in user.To")
	}

	tx = ar.db.Model(&to).
		Updates(user.To{Message: input.Message})

	if tx.RowsAffected == 0 {
		log.Error("no rows affected on update to message")
		return errors.New("data is up to date")
	}

	if tx.Error != nil {
		log.Error("error on update to message")
		return tx.Error
	}

	actionType := ""
	switch submission.Status {
	case "waiting":
		actionType = "approve"
	case "revised":
		actionType = "revise"
	case "rejected":
		actionType = "reject"
	}

	tx = ar.db.Model(&user.To{}).
		Joins("JOIN users ON tos.user_id = users.id").
		Where("tos.user_id = ? AND tos.submission_id = ?", userID, submission.ID).
		Update("action_type", actionType)

	if tx.Error != nil {
		log.Error("error on update action_type in 'to' table")
		return tx.Error
	}

	return nil
}

// SelectSubmissionById implements approve.Repository
func (ar *approverModel) SelectSubmissionById(userID string, id int) (approve.Core, error) {
	var dbsub uMod.Submission
	var toDetail admin.Users
	var ccDetail admin.Users
	var toDetails []admin.Users
	var ccDetails []admin.Users

	query := ar.db.
		Table("submissions").
		// Preload("Position").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN users ON tos.user_id = users.id").
		Joins("JOIN types ON submissions.type_id = types.id").
		Joins("JOIN positions ON positions.id = users.position_id").
		Where("users.id = ? AND submissions.id = ?", userID, id).
		Preload("Type").
		Preload("User").
		Preload("Tos", "submission_id = ?", id).
		Preload("Ccs", "submission_id = ?", id).
		Preload("Signs", "submission_id = ?", id).
		Find(&dbsub)

	for _, to := range dbsub.Tos {
		if err := ar.db.Where("id = ?", to.UserID).Preload("Position").Find(&toDetail).Error; err != nil {
			log.Error(err)
			return approve.Core{}, err
		}
		toDetails = append(toDetails, toDetail)
	}
	for _, cc := range dbsub.Ccs {
		if err := ar.db.Where("id = ?", cc.UserID).Preload("Position").Find(&ccDetail).Error; err != nil {
			log.Error(err)
			return approve.Core{}, err
		}
		ccDetails = append(ccDetails, ccDetail)
	}

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return approve.Core{}, errors.New("submission not found")
		}
		log.Error("failed to find submission:", query.Error.Error())
		return approve.Core{}, errors.New("failed to retrieve submission")
	}

	return handler.SubmissionToCore(toDetails, ccDetails, dbsub), nil
}

// SelectSubmissionApprove implements approve.Repository
func (ar *approverModel) SelectSubmissionAprrove(userID string, limit, offset int, search string) ([]approve.Core, error) {
	var res []approve.Core
	var dbsub []uMod.Submission

	query := ar.db.Table("submissions").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN types ON submissions.type_id = types.id").
		Where("tos.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Preload("Type").
		Find(&dbsub)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return []approve.Core{}, errors.New("submission not found")
		}
		log.Error("failed to find all submission:", query.Error.Error())
		return []approve.Core{}, errors.New("failed to retrieve all submission")
	}

	for _, v := range dbsub {
		tmp := approve.Core{
			ID:        v.ID,
			UserID:    v.UserID,
			TypeID:    v.TypeID,
			Title:     v.Title,
			Message:   v.Message,
			Status:    v.Status,
			Is_Opened: false,
			CreatedAt: time.Time{},
			Type: subtype.Core{
				SubmissionTypeName: v.Type.Name,
			},
		}
		res = append(res, tmp)
	}

	return res, nil
}
