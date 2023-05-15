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

func (sm *submissionModel) FindRequirement(userID string, typeName string, value string) (submission.RequirementDB, error) {
	var applicant admin.Users
	var typeDetail admin.Type
	if err := sm.db.Where("id = ?", userID).First(&applicant).Error; err != nil {
		log.Errorf("error on finding applicant with this id : %s, \n%w", userID, err)
		return submission.RequirementDB{}, err
	}

	var positions []admin.Position
	err := sm.db.Raw("SELECT positions.* FROM positions JOIN types ON positions.type_id = types.id WHERE types.name = ?", typeName).Scan(&positions).Error
	if err != nil {
		log.Errorf("error on finding positions with this type name : %s, \n%w", typeName, err)
		return submission.RequirementDB{}, err
	}

	var tos []admin.Users
	if err := sm.db.Joins("PositionHasType").Joins("Positions").Where("position_has_types.as = ? AND position_has_types.type_id IN ? AND position_has_types.to_level IS NOT NULL", "To", positions).Order("position_has_types.to_level ASC").Find(&tos).Error; err != nil {
		log.Errorf("error on finding positions with to role : %s, \n%w", typeName, err)
		return submission.RequirementDB{}, err
	}

	sm.db.Joins("PositionHasType").Joins("Positions").Where("position_has_types.as = ? AND position_has_types.type_id IN ? AND position_has_types.to_level IS NOT NULL", "To", positions).Order("position_has_types.to_level ASC").Find(&tos)

	var ccs []admin.Users
	if err := sm.db.Joins("PositionHasType").Joins("Positions").Where("position_has_types.as = ? AND position_has_types.type_id IN ? AND position_has_types.to_level IS NULL", "Cc", positions).Order("position_has_types.to_level ASC").Find(&ccs).Error; err != nil {
		log.Errorf("error on finding positions with cc role : %s, \n%w", typeName, err)
		return submission.RequirementDB{}, err
	}

	sm.db.Joins("Office").Where("users.office_id = ? OR offices.parent_id ?", applicant.OfficeID, applicant.Office.ParentID).Find(&tos)
	sm.db.Joins("Office").Where("users.office_id = ? OR offices.parent_id ?", applicant.OfficeID, applicant.Office.ParentID).Find(&ccs)

	if err := sm.db.Where("name = ?", typeName).First(&typeDetail).Error; err != nil {
		log.Errorf("errpr on finding typeDetail with name %s %w ", typeName, err)
		return submission.RequirementDB{}, err
	}

	var requirementDB submission.RequirementDB
	requirementDB.ApplicantName = applicant.Name
	requirementDB.ApplicantID = applicant.ID

	return submission.RequirementDB{}, nil
}
