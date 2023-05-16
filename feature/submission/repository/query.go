package repository

import (
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

	if err := sm.db.Raw(`
    SELECT users.id, users.name, users.email, users.phone_number, offices.id as office_id, position_has_types.as AS position_role
    FROM users
    INNER JOIN position_has_types ON position_has_types.position_id = users.position_id
    INNER JOIN positions ON positions.id = position_has_types.position_id
    INNER JOIN offices ON offices.id = users.office_id
    INNER JOIN types ON types.id = position_has_types.type_id
    WHERE positions.deleted_at IS NULL AND types.deleted_at IS NULL AND position_has_types.deleted_at IS NULL
    AND types.name = ? AND position_has_types.as = 'cc' AND position_has_types.value = ? AND (users.office_id = ? OR offices.name = 'Head Office')
`, typeName, 30000000, applicant.Office.ID).Scan(&ccs).Error; err != nil {
		return submission.Core{}, err
	}

	if err := sm.db.Raw(`
    SELECT users.id, users.name, users.email, users.phone_number, offices.id as office_id, position_has_types.as AS position_role
    FROM users
    INNER JOIN position_has_types ON position_has_types.position_id = users.position_id
    INNER JOIN positions ON positions.id = position_has_types.position_id
    INNER JOIN offices ON offices.id = users.office_id
    INNER JOIN types ON types.id = position_has_types.type_id
    WHERE positions.deleted_at IS NULL AND types.deleted_at IS NULL AND position_has_types.deleted_at IS NULL
    AND types.name = ? AND position_has_types.as = 'to' AND position_has_types.value = ? AND (users.office_id = ? OR offices.name = 'Head Office')
	ORDER BY position_has_types.to_level ASC
`, typeName, 30000000, applicant.Office.ID).Scan(&tos).Error; err != nil {
		return submission.Core{}, err
	}

	var result submission.Core

	result.Requirement = typeDetail.Requirement

	for _, to := range tos {
		tmp := submission.ToApprover{
			ApproverPosition: to.Position.Name,
			ApproverId:       to.ID,
			ApproverName:     to.Name,
		}
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

	for _, v := range  newSub.ToApprover{
		tmp := To{
			ID: v.ApproverId,
			Name: v.ApproverName,
			
		}
		
	}


	err := sm.db.
}
