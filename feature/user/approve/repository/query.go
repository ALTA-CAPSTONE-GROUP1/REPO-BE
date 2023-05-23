package repository

import (
	"errors"
	"fmt"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"

	// uMod "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve/handler"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
)

type approverModel struct {
	db *gorm.DB
	u  helper.UploadInterface
}

func New(db *gorm.DB, u helper.UploadInterface) approve.Repository {
	return &approverModel{
		db: db,
		u:  u,
	}
}

// UpdateUser implements approve.Repository
func (ar *approverModel) UpdateApprove(userID string, id int, input approve.Core) error {
	var submission user.Submission
	var tos []user.To
	var owner admin.Users

	tx := ar.db.Model(&submission).
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN users ON tos.user_id = users.id").
		Where("tos.user_id = ? AND submissions.id = ?", userID, id).
		Find(&submission)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}
	fmt.Println(submission)

	tx = ar.db.Where("id = ?", submission.UserID).First(&owner)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}
	fmt.Println(owner)

	tx = ar.db.Where("submission_id = ?", submission.ID).Find(&tos)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}
	fmt.Println(tos)

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

	if len(tos) > 0 && tos[len(tos)-1].UserID == userID {
		submission.Status = "Approved"
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
	switch submission.Status {
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
	var signdb user.Sign
	signdb.UserID = userID
	signdb.Name = sign
	signdb.SubmissionID = submission.ID
	tx = ar.db.Create(&signdb)

	if tx.Error != nil {
		log.Error("error on update sign")
		return tx.Error
	}

	file := user.File{}
	tx = ar.db.Model(&user.File{}).
		Select("link").
		Where("submission_id = ?", submission.ID).
		First(&file)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}

	// approverProfile := admin.Users{}
	// tx = ar.db.Where("id = ?", userID).Find(&approverProfile)
	// if tx.Error != nil {
	// 	log.Errorf("error on finding owner%w", tx.Error)
	// 	return tx.Error
	// }

	recipient := []string{owner.Email}
	receiverName := []string{owner.Name}
	helper.SendSimpleEmail(signdb.Name, submission.Title, "Update on your submission", recipient, receiverName, owner.Email)
	// currentLink := file.Link
	// currentFileName := file.Name
	// submissionTitle := submission.Title
	// signName := sign
	// action := actionType
	// actionMessage := input.Message
	// approverName := to.User.Name
	// approverPosition := to.User.Position.Name

	// fileHeader, err := helper.UpdateCreateSign(currentLink, currentFileName, approverName, submissionTitle, signName, approverPosition, action, actionMessage)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(fileHeader)

	// newpath := helper.GenerateIDFromPositionTag(userID)
	// realpath := "/" + newpath + "/" + userID
	// newLink, err := ar.u.UploadFile(fileHeader, realpath)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(newLink)

	// var newFile user.File
	// newFile.Link = newLink[0]
	// newFile.Name = (fileHeader.Filename + newpath)
	// newFile.SubmissionID = submission.ID

	// tx = ar.db.Where("submission_id = ?, id = ?", submission.ID, file.ID).Update("link", newFile.Link)
	// time.Sleep(3 * time.Second)
	// if tx.Error != nil {
	// 	log.Errorf("errros on updating file lin %w", tx.Error)
	// 	return tx.Error
	// }

	return nil
}

// SelectSubmissionById implements approve.Repository
func (ar *approverModel) SelectSubmissionById(userID string, id int) (approve.Core, error) {
	var dbsub user.Submission
	var toDetails []admin.Users
	var ccDetails []admin.Users
	var fileDetails []user.File

	if err := ar.db.Model(&user.To{}).Where("id = ?", id).Update("is_opened", 1).Error; err != nil {
		log.Errorf("error in update status opened", err)
		return approve.Core{}, err
	}
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

	var fileDetail user.File
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
	var dbsub []user.Submission

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
		Preload("Tos").
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
			ID:     v.ID,
			UserID: v.UserID,
			Owner: approve.OwnerCore{
				Name: v.User.Name,
			},
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
