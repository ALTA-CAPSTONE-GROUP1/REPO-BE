package repository

import (
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
		Where("types.name = ? AND position_has_types.as = 'Cc' AND position_has_types.value = ? AND (users.office_id = ? OR offices.name = 'Head Office')", typeName, 30000000, applicant.Office.ID).
		Find(&ccs).Error; err != nil {
		return submission.Core{}, err
	}

	if err := sm.db.Preload("Position").
		Joins("INNER JOIN position_has_types ON position_has_types.position_id = users.position_id").
		Joins("INNER JOIN positions ON positions.id = position_has_types.position_id").
		Joins("INNER JOIN offices ON offices.id = users.office_id").
		Joins("INNER JOIN types ON types.id = position_has_types.type_id").
		Where("positions.deleted_at IS NULL AND types.deleted_at IS NULL AND position_has_types.deleted_at IS NULL").
		Where("types.name = ? AND position_has_types.as = 'To' AND position_has_types.value = ? AND (users.office_id = ? OR offices.name = 'Head Office')", typeName, 30000000, applicant.Office.ID).
		Order("position_has_types.to_level ASC").
		Find(&tos).Error; err != nil {
		return submission.Core{}, err
	}
	// fmt.Println(tos)
	fmt.Println(tos[0].Position)
	var result submission.Core

	result.Requirement = typeDetail.Requirement

	for _, to := range tos {
		tmp := submission.ToApprover{
			ApproverPosition: to.Position.Name,
			ApproverId:       to.ID,
			ApproverName:     to.Name,
		}
		fmt.Println(to.Position.Name)
		result.To = append(result.To, tmp)
	}

	for _, cc := range ccs {
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

	for _, v := range newSub.ToApprover {
		tmp := To{
			Name:   v.ApproverName,
			UserID: v.ApproverId,
		}
		submissionDB.Tos = append(submissionDB.Tos, tmp)
	}

	for _, v := range newSub.CC {
		tmp := Cc{
			Name:   v.CcName,
			UserID: v.CcId,
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

	for _, v := range user.Position.Types {
		var phts []admin.PositionHasType
		if err := sm.db.Where("position_id = ? AND type_id = ? AND `as` = ?", user.Position.ID, v.ID, "Owner").Find(&phts).Error; err != nil {
			log.Errorf("error on finding subTypes have by user", err)
			return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, err
		}
		existingIndex := -1
		for i, choice := range choices {
			if choice.SubTypeName == v.Name {
				existingIndex = i
				break
			}
		}

		if existingIndex != -1 {
			for _, detail := range phts {
				choices[existingIndex].SubtypeValue = append(choices[existingIndex].SubtypeValue, detail.Value)
			}
		} else {
			subTypeChoices := submission.SubTypeChoices{
				SubTypeName:  v.Name,
				SubtypeValue: make([]int, 0, len(phts)),
			}
			for _, detail := range phts {
				subTypeChoices.SubtypeValue = append(subTypeChoices.SubtypeValue, detail.Value)
			}
			choices = append(choices, subTypeChoices)
		}
	}

	fmt.Println(choices)

	if err := sm.db.Where("user_id = ?", userID).
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
			if err := sm.db.Where("id = ?", to.UserID).Preload("Position").First(&toDetails).Error; err != nil {
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
			if err := sm.db.Where("id = ?", cc.UserID).Preload("Position").First(&ccDetails).Error; err != nil {
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
			ReceiveDate:    sub.CreatedAt.Format(time.RFC3339),
			Opened:         sub.Is_Opened,
			Attachment:     attachment.Link,
			SubmissionType: subTypeByID.Name,
		})
	}

	return resultAllSubmission, choices, nil
}

func (sm *submissionModel) SelectSubmissionByID(submissionID int) (submission.GetSubmissionByIDCore, error) {
	var result submission.GetSubmissionByIDCore

	var submissions Submission
	if err := sm.db.Where("id = ?", submissionID).
		Preload("Files").
		Preload("Tos").
		Preload("Ccs").
		First(&submissions).Error; err != nil {
		log.Errorf("error on finding all data from submission ID: %v", err)
		return submission.GetSubmissionByIDCore{}, fmt.Errorf("error on finding submissionDatas: %w", err)
	}

	result.Attachment = submissions.Files[0].Link

	var typeDetails admin.Type
	if err := sm.db.Where("id = ?", submissions.TypeID).
		First(&typeDetails).Error; err != nil {
		log.Errorf("error on finding submissionType detail %w", err)
		return submission.GetSubmissionByIDCore{},
			fmt.Errorf("error on finding submissionType %w", err)
	}

	result.SubmissionType = typeDetails.Name

	for _, submissionTo := range submissions.Tos {
		tmp := submission.ToApprover{
			ApproverPosition: submissionTo.Name,
			ApproverId:       submissionTo.UserID,
		}
		tmpaction := submission.ApproverActions{
			Action:       submissionTo.Action_Type,
			ApproverName: submissionTo.Name,
			Message:      submissionTo.Message,
		}
		result.To = append(result.To, tmp)
		result.ApproverActions = append(result.ApproverActions, tmpaction)
	}

	for i, v := range result.To {
		var to admin.Users
		if err := sm.db.Preload("Positions").
			Where("id = ?", v.ApproverId).
			First(&to).Error; err != nil {
			log.Errorf("eror on finding to position %w", err)
			return submission.GetSubmissionByIDCore{},
				fmt.Errorf("error on to positions %w", err)
		}
		result.To[i].ApproverPosition = to.Position.Name
		result.ApproverActions[i].ApproverPosition = to.Position.Name
	}

	for _, submissionCc := range submissions.Ccs {
		tmp := submission.CcApprover{
			CcName: submissionCc.Name,
			CcId:   submissionCc.UserID,
		}

		result.CC = append(result.CC, tmp)
	}

	for i, v := range result.CC {
		var cc admin.Users
		if err := sm.db.Preload("Positions").
			Where("id = ?", v.CcId).
			First(&cc).Error; err != nil {
			log.Errorf("eror on finding to position %w", err)
			return submission.GetSubmissionByIDCore{},
				fmt.Errorf("error on to positions %w", err)
		}

		result.CC[i].CcPosition = cc.Position.Name
	}

	result.Message = submissions.Message
	result.Title = submissions.Title

	return result, nil
}
