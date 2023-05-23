package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type submissionModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) submission.Repository {
	return &submissionModel{
		db: db,
	}
}

func (sm *submissionModel) FindRequirement(userID string, typeName string, typeValue int) (submission.Core, error) {
	var applicant admin.Users
	var typeDetail admin.Type
	var tos []admin.Users
	var ccs []admin.Users
	applicant.ID = userID

	if err := sm.db.Where("name = ? ", typeName).Find(&typeDetail).Error; err != nil {
		log.Errorf("error on finding typeDetails from %s", typeName)
		return submission.Core{}, err
	}

	if err := sm.db.Preload("Office").First(&applicant).Error; err != nil {
		log.Errorf("error on finding applicant office %s", applicant.Name)
		return submission.Core{}, err
	}

	if err := sm.db.Preload("Position").
		Joins("INNER JOIN position_has_types ON position_has_types.position_id = users.position_id").
		Joins("INNER JOIN positions ON positions.id = position_has_types.position_id").
		Joins("INNER JOIN offices ON offices.id = users.office_id").
		Joins("INNER JOIN types ON types.id = position_has_types.type_id").
		Where("positions.deleted_at IS NULL AND types.deleted_at IS NULL AND position_has_types.deleted_at IS NULL").
		Where("types.name = ? AND position_has_types.as = 'Cc' AND position_has_types.value = ? AND (users.office_id = ? OR offices.name = 'Head Office')", typeName, typeValue, applicant.Office.ID).
		Find(&ccs).Error; err != nil {
		return submission.Core{}, err
	}
	var findTypeId int

	if err := sm.db.Raw("SELECT id FROM types WHERE name = ?", typeName).Find(&findTypeId).Error; err != nil {
		return submission.Core{}, err
	}

	fmt.Println(findTypeId)
	fmt.Println(applicant.Position.ID)

	var applicantToLevel int
	if err := sm.db.Raw("SELECT to_level FROM position_has_types WHERE type_id = ? AND position_id = ?", &findTypeId, applicant.Position.ID).Find(&applicantToLevel).Error; err != nil {
		return submission.Core{}, err
	}

	if err := sm.db.Preload("Position").
		Joins("INNER JOIN position_has_types ON position_has_types.position_id = users.position_id").
		Joins("INNER JOIN positions ON positions.id = position_has_types.position_id").
		Joins("INNER JOIN offices ON offices.id = users.office_id").
		Joins("INNER JOIN types ON types.id = position_has_types.type_id").
		Where("positions.deleted_at IS NULL AND types.deleted_at IS NULL AND position_has_types.deleted_at IS NULL").
		Where("types.name = ? AND position_has_types.as = 'To' AND position_has_types.value = ? AND (users.office_id = ? OR offices.name = 'Head Office')", typeName, typeValue, applicant.Office.ID).
		Order("position_has_types.to_level ASC").
		Find(&tos).Error; err != nil {
		return submission.Core{}, err
	}

	var result submission.Core

	result.Requirement = typeDetail.Requirement

	for _, to := range tos {
		if to.ID == userID {
			continue
		}

		tmp := submission.ToApprover{
			ApproverPosition: to.Position.Name,
			ApproverId:       to.ID,
			ApproverName:     to.Name,
		}
		result.To = append(result.To, tmp)
	}

	for _, cc := range ccs {
		if cc.ID == userID {
			continue
		}

		tmp := submission.CcApprover{
			CcPosition: cc.Position.Name,
			CcId:       cc.ID,
			CcName:     cc.Name,
		}
		result.CC = append(result.CC, tmp)
	}

	return result, nil
}

func (sm *submissionModel) InsertSubmission(newSub submission.AddSubmissionCore) error {
	var submissionDB Submission

	var SubmissionTypeQuery admin.Type

	if err := sm.db.Where("name = ?", newSub.SubmissionType).First(&SubmissionTypeQuery).Error; err != nil {
		log.Error("cannot find submissiontype by name")
		return err
	}

	submissionDB.UserID = newSub.OwnerID
	submissionDB.Title = newSub.Title
	submissionDB.TypeID = SubmissionTypeQuery.ID
	submissionDB.Is_Opened = false
	submissionDB.UserID = newSub.OwnerID
	submissionDB.Status = "Sent"
	submissionDB.Message = newSub.Message

	for _, v := range newSub.ToApprover {
		tmp := To{
			Name:   v.ApproverName,
			UserID: v.ApproverId,
		}
		submissionDB.Tos = append(submissionDB.Tos, tmp)
	}

	for _, x := range newSub.CC {
		tmp := Cc{
			Name:   x.CcName,
			UserID: x.CcId,
		}
		submissionDB.Ccs = append(submissionDB.Ccs, tmp)
	}

	file := File{
		Name: newSub.Attachment,
		Link: newSub.AttachmentLink,
	}
	submissionDB.Files = append(submissionDB.Files, file)

	if err := sm.db.Create(&submissionDB).Error; err != nil {
		log.Error("error occurs while insert submission datas")
		return err
	}

	return nil
}

func (sm *submissionModel) SelectAllSubmissions(userID string, pr submission.GetAllQueryParams) ([]submission.AllSubmiisionCore, []submission.SubTypeChoices, error) {
	var (
		dbsubmissions       []Submission
		resultAllSubmission []submission.AllSubmiisionCore
		user                admin.Users
		choices             []submission.SubTypeChoices
	)

	if err := sm.db.Where("id = ?", userID).Preload("Position.Types").Find(&user).Error; err != nil {
		log.Errorf("error on finding subTypes have by user", err)
		return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, err
	}
	var phts []admin.PositionHasType
	var phtss []admin.PositionHasType
	for _, v := range user.Position.Types {

		if err := sm.db.Where("position_id = ? AND type_id = ? AND `as` = ?", user.Position.ID, v.ID, "Owner").Preload("Type").Find(&phts).Error; err != nil {
			log.Errorf("error on finding subTypes have by user", err)
			return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, err
		}

		phtss = append(phtss, phts...)
	}

	typeValues := make(map[string][]int)
	valueExist := make(map[string]map[int]bool)

	for _, v := range phtss {
		key := v.Type.Name
		value := v.Value

		if valueExist[key] == nil {
			typeValues[key] = append(typeValues[key], value)
			valueExist[key] = map[int]bool{value: true}
		} else {
			if !valueExist[key][value] {
				typeValues[key] = append(typeValues[key], value)
				valueExist[key][value] = true
			}
		}
	}

	for key, v := range typeValues {
		tmp := submission.SubTypeChoices{
			SubTypeName:  key,
			SubtypeValue: v,
		}
		choices = append(choices, tmp)
	}

	if err := sm.db.Where("user_id = ?", userID).Order("created_at DESC").
		Preload("Files").
		Preload("Tos").
		Preload("Ccs").
		Preload("Signs").
		Find(&dbsubmissions).Error; err != nil {
		log.Errorf("error on finding submissions for user %s: %v", userID, err)
		return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, err
	}

	for _, sub := range dbsubmissions {
		var toApprover []submission.ToApprover
		for _, to := range sub.Tos {
			var toDetails admin.Users
			if err := sm.db.Where("id = ?", to.UserID).Preload("Position").Find(&toDetails).Error; err != nil {
				log.Errorf("failed on finding positions of tos %w", err)
				return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, err
			}
			toApprover = append(toApprover, submission.ToApprover{
				ApproverId:       to.UserID,
				ApproverName:     toDetails.Name,
				ApproverPosition: toDetails.Position.Name,
			})
		}
		var ccApprover []submission.CcApprover
		for _, cc := range sub.Ccs {
			var ccDetails admin.Users
			if err := sm.db.Where("id = ?", cc.UserID).Preload("Position").Find(&ccDetails).Error; err != nil {
				log.Error("failed on finding positions of tos")
				return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, err
			}
			ccApprover = append(ccApprover, submission.CcApprover{
				CcPosition: ccDetails.Position.Name,
				CcName:     ccDetails.Name,
				CcId:       cc.UserID,
			})
		}

		var subTypeByID admin.Type
		if err := sm.db.Where("id = ?", sub.TypeID).First(&subTypeByID).Error; err != nil {
			log.Errorf("error getting files for subType %d: %v", sub.TypeID, err)
			return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, err
		}

		var attachment File
		if err := sm.db.Where("submission_id = ?", sub.ID).First(&attachment).Error; err != nil {
			log.Errorf("error getting files for submission %d: %v", sub.ID, err)
			return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, err
		}

		resultAllSubmission = append(resultAllSubmission, submission.AllSubmiisionCore{
			ID:             sub.ID,
			Tos:            toApprover,
			CCs:            ccApprover,
			Title:          sub.Title,
			Status:         sub.Status,
			ReceiveDate:    sub.CreatedAt.Add(7 * time.Hour).Format("2006-01-02 15:04"),
			Opened:         sub.Is_Opened,
			Attachment:     attachment.Link,
			SubmissionType: subTypeByID.Name,
		})
	}

	return resultAllSubmission, choices, nil
}

func (sm *submissionModel) SelectSubmissionByID(submissionID int, userID string) (submission.GetSubmissionByIDCore, error) {
	var (
		result         submission.GetSubmissionByIDCore
		submissionByID Submission
		subTypeDetails admin.Type
	)
	if err := sm.db.Model(&Submission{}).Where("id = ?", submissionID).Update("is_opened", 1).Error; err != nil {
		log.Errorf("error in update status opened", err)
		return submission.GetSubmissionByIDCore{}, err
	}

	if err := sm.db.Where("user_id = ? AND id = ?", userID, submissionID).
		Preload("Files").
		Preload("Tos").
		Preload("Ccs").
		Preload("Signs").
		First(&submissionByID).Error; err != nil {
		log.Errorf("error on finding submissions for by userid and submissionid %s: %v", userID, err)
		return submission.GetSubmissionByIDCore{}, err
	}

	if err := sm.db.Where("id = ?", submissionByID.TypeID).Find(&subTypeDetails).Error; err != nil {
		log.Errorf("error in finding submissionType Name %w", err)
		return submission.GetSubmissionByIDCore{}, err
	}

	var toApprover []submission.ToApprover
	var toActions []submission.ApproverActions
	for _, to := range submissionByID.Tos {
		var toDetails admin.Users
		if err := sm.db.Where("id = ?", to.UserID).Preload("Position").First(&toDetails).Error; err != nil {
			log.Errorf("failed on finding positions of tos %w", err)
			return submission.GetSubmissionByIDCore{}, err
		}
		toApprover = append(toApprover, submission.ToApprover{
			ApproverId:       to.UserID,
			ApproverName:     toDetails.Name,
			ApproverPosition: toDetails.Position.Name,
		})
		toActions = append(toActions, submission.ApproverActions{
			Action:           to.Action_Type,
			ApproverName:     toDetails.Name,
			ApproverPosition: toDetails.Position.Name,
			Message:          to.Message,
		})
		var signs []Sign
		if err := sm.db.Where("id = ? AND user_id = ?", submissionByID.ID, to.UserID).Find(&signs).Error; err != nil {
			log.Errorf("error on finding sign datas %w", err)
			return submission.GetSubmissionByIDCore{}, err
		}

		for _, v := range signs {
			result.Signs = append(result.Signs, v.Name)
		}
	}
	var ccApprover []submission.CcApprover
	for _, cc := range submissionByID.Ccs {
		var ccDetails admin.Users
		if err := sm.db.Where("id = ?", cc.UserID).Preload("Position").First(&ccDetails).Error; err != nil {
			log.Errorf("failed on finding positions of tos %w", err)
			return submission.GetSubmissionByIDCore{}, err
		}
		ccApprover = append(ccApprover, submission.CcApprover{
			CcPosition: ccDetails.Position.Name,
			CcName:     ccDetails.Name,
			CcId:       cc.UserID,
		})
	}

	if len(submissionByID.Files) > 0 {
		result.Attachment = submissionByID.Files[0].Link
	}

	result.To = append(result.To, toApprover...)
	result.ApproverActions = append(result.ApproverActions, toActions...)
	result.CC = append(result.CC, ccApprover...)
	result.Title = submissionByID.Title
	result.Message = submissionByID.Message
	result.Status = submissionByID.Status
	if len(toActions) > 0 {
		result.ActionMessage = toActions[len(toActions)-1].Message
	}
	result.ActionMessage = toActions[0].Message

	result.SubmissionType = subTypeDetails.Name

	return result, nil
}

func (sm *submissionModel) DeleteSubmissionByID(submissionID int, userID string) error {
	var submission Submission
	var count int64
	if err := sm.db.Where("id = ? AND user_id = ?", submissionID, userID).Find(&submission).Count(&count).Error; err != nil {
		log.Warn("cannot find submission datas")
		return err
	}
	if count == 0 {
		return fmt.Errorf("submission data not found")
	}
	if submission.Status != "Sent" {
		log.Warn("submission status not 'sent'")
		return errors.New("submission status are not 'Sent'")
	}
	if err := sm.db.Delete(&submission).Error; err != nil {
		log.Errorf("failed to delete submission")
		return err
	}

	return nil
}

func (sm *submissionModel) UpdateDataByOwner(editedData submission.UpdateCore) error {
	tx := sm.db.Begin()

	var submission Submission
	if err := tx.Where("id = ?", editedData.SubmissionID).First(&submission).Error; err != nil {
		tx.Rollback()
		log.Error("cannot find submission data")
		return errors.New("submission data not found")
	}
	if submission.Status != "Sent" && submission.Status != "Revised" {
		tx.Rollback()
		log.Errorf("submisison status not 'sent'")
		return errors.New("submission status not 'sent'")
	}

	if err := tx.Exec(`
		UPDATE submissions
		SET message = ?, status = ?
		WHERE id = ? AND user_id = ?
	`, editedData.Message, "Sent", editedData.SubmissionID, editedData.UserID).Error; err != nil {
		tx.Rollback()
		log.Errorf("error on update submission data")
		return err
	}

	if err := tx.Exec(`
		UPDATE files
		SET name = ?, link = ?
		WHERE submission_id = ?
	`, editedData.AttachmentName, editedData.AttachmentLink, editedData.SubmissionID).Error; err != nil {
		tx.Rollback()
		log.Errorf("error on saving new attachment")
		return err
	}

	if err := tx.Exec(`
		DELETE FROM signs 
		WHERE submission_id = ?
	`, editedData.SubmissionID).Error; err != nil {
		tx.Rollback()
		log.Errorf("error on delete signs data in database")
		return err
	}

	tx.Commit()
	return nil
}

func (sm *submissionModel) FindFileData(submissionID int, fileName string) bool {
	var file File
	result := sm.db.Where("name = ? AND submission_id = ?", fileName, submissionID).First(&file)

	if result.Error != nil {
		log.Errorf("error querying file data %s", result.Error)
	}
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false
	}

	return true
}
