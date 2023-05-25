package repository

import (
	"errors"
	"fmt"
	"time"

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

	tx = ar.db.Where("id = ?", submission.UserID).First(&owner)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}

	tx = ar.db.Where("submission_id = ?", submission.ID).Find(&tos)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}

	switch input.Status {
	case "approve":
		submission.Status = "Waiting"
		if len(tos) > 0 && tos[len(tos)-1].UserID == userID {
			submission.Status = "Approved"
		}
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

	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}

	tx = ar.db.Model(&to).
		Updates(user.To{Message: input.Message})

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
		Select("id, link").
		Where("submission_id = ?", submission.ID).
		First(&file)
	if tx.Error != nil {
		log.Errorf("error on finding owner%w", tx.Error)
		return tx.Error
	}

	recipient := []string{owner.Email}
	receiverName := []string{owner.Name}
	helper.SendSimpleEmail(signdb.Name, submission.Title, "Update on your submission", recipient, receiverName, owner.Email)

	action := input.Status
	curentLink := file.Link
	subTitle := submission.Title
	signName := sign
	newApproverName := to.UserID
	approverPosition := to.User.Position.Name
	newPath := helper.GenerateIDFromPositionTag(userID)
	realPath := "/" + newPath + time.Now().Format(time.RFC3339)

	newFileName, newLink, err := helper.UpdateFile(action, curentLink, newApproverName, approverPosition, subTitle, signName, realPath)
	if err != nil {
		log.Errorf("error on creating and uploading pdf stamp %w", err)
		return err
	}

	fmt.Print(newLink)
	if err := ar.db.Statement.Exec("SET FOREIGN_KEY_CHECKS=?", 0).Error; err != nil {
		log.Error(err.Error())
		return err
	}

	// saveFile := &user.File{
	// 	SubmissionID: submission.ID,
	// 	Link:         newLink,
	// 	Name:         newFileName,
	// }

	updateObj := user.File{
		ID: file.ID,
	}
	updates := map[string]interface{}{
		"Name": newFileName,
		"Link": newLink,
	}

	fmt.Printf("FILEID %d\n", file.ID)
	fmt.Printf("submissionid %d\n", submission.ID)
	if err := ar.db.Model(&updateObj).Where("id = ? AND submission_id = ?", file.ID, submission.ID).Updates(updates).Error; err != nil {
		log.Errorf("error on updating file: %s", err.Error())
		return err
	}

	if tx.Error != nil {
		log.Errorf("error on saving file %s", err.Error())
		return err
	}

	if tx.RowsAffected == 0 {
		log.Warn("no rows affected in updaing file")
	}

	if err := ar.db.Statement.Exec("SET FOREIGN_KEY_CHECKS=?", 0).Error; err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// SelectSubmissionById implements approve.Repository
func (ar *approverModel) SelectSubmissionById(userID string, id int) (approve.Core, error) {
	var dbsub user.Submission
	var toDetails []admin.Users
	var ccDetails []admin.Users
	var fileDetails []user.File

	if err := ar.db.Model(&user.To{}).Where("submission_id = ? AND user_id = ?", id, userID).
		Update("is_opened", 1).Error; err != nil {
		log.Errorf("error in update status opened: %v", err)
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
func (ar *approverModel) SelectSubmissionAprrove(userID string, search approve.GetAllQueryParams) ([]approve.Core, int, error) {
	var res []approve.Core
	var dbsub []user.Submission
	totalData := int64(-1)

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
		Preload("Tos", "tos.user_id = ?", userID).
		Preload("User.Position").
		Where("tos.user_id = ?", userID).
		Order("created_at DESC")

	if title != "" {
		query = query.Where("submissions.title LIKE ?", "%"+title+"%")
	}

	if fromParm != "" {
		query = query.Where("positions.name LIKE ?", "%"+fromParm+"%")
	}

	if types != "" {
		query = query.Where("types.name LIKE ?", "%"+types+"%")
	}

	query.Model(&user.Submission{}).Count(&totalData)

	query = query.Limit(limit).
		Offset(offset).
		Preload("Type").
		Find(&dbsub)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return []approve.Core{}, 0, errors.New("submission not found")
		}
		log.Error("failed to find all submissions:", query.Error.Error())
		return []approve.Core{}, 0, errors.New("failed to retrieve all submissions")
	}

	for _, v := range dbsub {
		tmp := approve.Core{
			ID:     v.ID,
			UserID: v.UserID,
			Owner: approve.OwnerCore{
				Name:     v.User.Name,
				Position: v.User.Position.Name,
			},
			TypeID:    v.TypeID,
			Title:     v.Title,
			Message:   v.Message,
			Is_Opened: false,
			CreatedAt: v.CreatedAt,
			Type: subtype.Core{
				SubmissionTypeName: v.Type.Name,
			},
		}

		for _, to := range v.Tos {
			if to.UserID == userID {
				actionType := to.Action_Type
				if v.Status == "Sent" {
					actionType = "waiting for you"

					ar.db.Model(&to).Update("Action_Type", actionType)
					ar.db.Save(&to)

				}
				tmp.Tos = append(tmp.Tos, approve.ToCore{
					SubmissionID: to.SubmissionID,
					UserID:       to.UserID,
					Name:         to.Name,
					Action_Type:  actionType,
				})
			}
		}

		res = append(res, tmp)
	}

	return res, int(totalData), nil
}
