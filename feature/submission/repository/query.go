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

func (sm *submissionModel) FindRequirement(applicantID string, typeName string, value string) (submission.RequirementDB, error) {
	var applicant admin.Users
	var typeDetail admin.Type
	if err := sm.db.Where("id = ?", applicantID).First(&applicant).Error; err != nil {
		log.Errorf("error on finding applicant with this id : %s, \n%w", applicantID, err)
		return submission.RequirementDB{}, err
	}
	//syntax error
	//record not found
	//unable to connect to database

	if err := sm.db.Where("name = ?", typeName).First(&typeDetail).Error; err != nil {
		log.Errorf("errpr on finding typeDetail with name %s %w ", typeName, err)
		return submission.RequirementDB{}, err
	}

	var requirementDB submission.RequirementDB
	requirementDB.ApplicantName = applicant.Name
	requirementDB.ApplicantID = applicant.ID

	return submission.RequirementDB{}, nil
}
