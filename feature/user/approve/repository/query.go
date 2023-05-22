package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	uMod "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve/handler"
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
		submission.Status = "Waiting"
	case "revise":
		submission.Status = "Revised"
	case "reject":
		submission.Status = "Rejected"
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
	case "Waiting":
		actionType = "approve"
	case "Revised":
		actionType = "revise"
	case "Rejected":
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

	sign, err := helper.GenerateUniqueSign(userID)
	if err != nil {
		log.Error("failed to generate unique sign")
		return err
	}

	tx = ar.db.Model(&user.Sign{}).
		Where("signs.user_id = ? AND signs.submission_id = ?", userID, submission.ID).
		Update("sign.name", sign)

	if tx.RowsAffected == 0 {
		log.Error("no rows affected on update sign")
		return errors.New("data is up to date")
	}

	if tx.Error != nil {
		log.Error("error on update sign")
		return tx.Error
	}

	// file := user.File{}
	// tx = ar.db.Model(&user.File{}).
	// 	Select("link").
	// 	Where("submission_id = ?", submission.ID).
	// 	First(&file)

	// if tx.RowsAffected == 0 {
	// 	log.Error("no file found for the given submission ID")
	// 	return errors.New("no file found")
	// }

	// currentLink := fmt.Sprintf("link/to/submission/%d", submission.ID)
	// currentFileName := file.Link
	// submissionTitle := submission.Title
	// signName := sign
	// action := input.Status
	// actionMessage := input.Message
	// approverName := to.User.Name
	// approverPosition := to.User.Position.Name

	// fileHeader, err := helper.UpdateCreateSign(currentLink, currentFileName, approverName, submissionTitle, signName, approverPosition, action, actionMessage)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// SelectSubmissionById implements approve.Repository
func (ar *approverModel) SelectSubmissionById(userID string, id int) (approve.Core, error) {
	var dbsub uMod.Submission
	var toDetails []admin.Users
	var ccDetails []admin.Users
	var fileDetails []uMod.File

	query := ar.db.
		Table("submissions").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN users ON tos.user_id = users.id").
		Joins("JOIN types ON submissions.type_id = types.id").
		Joins("JOIN positions ON positions.id = users.position_id").
		Where("users.id = ? AND submissions.id = ?", userID, id).
		Preload("Type").
		Preload("Tos").
		Preload("Ccs").
		Preload("Signs").
		Find(&dbsub)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return approve.Core{}, errors.New("submission not found")
		}
		log.Error("failed to find submission:", query.Error.Error())
		return approve.Core{}, errors.New("failed to retrieve submission")
	}

	for _, to := range dbsub.Tos {
		var toDetail admin.Users
		if err := ar.db.Where("id = ?", to.UserID).Preload("Position").Find(&toDetail).Error; err != nil {
			log.Error(err)
			return approve.Core{}, err
		}
		toDetails = append(toDetails, toDetail)
	}

	for _, cc := range dbsub.Ccs {
		var ccDetail admin.Users
		if err := ar.db.Where("id = ?", cc.UserID).Preload("Position").Find(&ccDetail).Error; err != nil {
			log.Error(err)
			return approve.Core{}, err
		}
		ccDetails = append(ccDetails, ccDetail)
	}

	var fileDetail uMod.File
	if err := ar.db.Where("submission_id = ?", dbsub.ID).Find(&fileDetail).Error; err != nil {
		log.Error(err)
		return approve.Core{}, err
	}
	fileDetails = append(fileDetails, fileDetail)

	var owner admin.Users
	if err := ar.db.Where("id = ?", dbsub.UserID).Preload("Position").Find(&owner).Error; err != nil {
		log.Error(err)
		return approve.Core{}, err
	}

	return handler.SubmissionToCore(owner, fileDetails, toDetails, ccDetails, dbsub), nil
}

// SelectSubmissionApprove implements approve.Repository
func (ar *approverModel) SelectSubmissionAprrove(userID string, search approve.GetAllQueryParams) ([]approve.Core, error) {
	var res []approve.Core
	var dbsub []uMod.Submission

	limit := search.Limit
	offset := search.Offset
	fromParm := search.FromTo
	title := search.Title
	types := search.Type

	query := ar.db.Table("submissions").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN users ON tos.user_id = users.id").
		Joins("JOIN types ON submissions.type_id = types.id").
		Joins("JOIN positions ON positions.id = users.position_id").
		Preload("User").
		Preload("Type").
		Where("tos.user_id = ?", userID)

	if title != "" {
		query = query.Where("submissions.title LIKE ?", "%"+title+"%")
	}

	if fromParm != "" {
		query = query.Where("positions.name LIKE ?", "%"+fromParm+"%")
	}

	if types != "" {
		query = query.Where("types.name LIKE ?", "%"+types+"%")
	}

	query = query.Limit(limit).
		Offset(offset).
		Preload("Type").
		Find(&dbsub)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return []approve.Core{}, errors.New("submission not found")
		}
		log.Error("failed to find all submissions:", query.Error.Error())
		return []approve.Core{}, errors.New("failed to retrieve all submissions")
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
			CreatedAt: v.CreatedAt,
			Type: subtype.Core{
				SubmissionTypeName: v.Type.Name,
			},
		}
		res = append(res, tmp)
	}

	return res, nil
}
