package repository

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
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

func (sm *submissionModel) FindRequirement(userID string, typeName string, typeValue int) (submission.RequirementDB, error) {
	// var applicant admin.Users
	// var typeDetail admin.Type
	// var tos []admin.Users

	// sm.db.Raw(`
	// 	SELECT users.*
	// 	FROM users
	// 	INNER JOIN positions ON users.position_id = positions.id
	// 	INNER JOIN position_has_types ON positions.id = position_has_types.position_id AND position_has_types.type_id = ? AND position_has_types.value = ? AND position_has_types.as = 'to'
	// 	INNER JOIN offices ON users.office_id = offices.id AND (offices.id = ? OR offices.parent_id = ?)
	// 	INNER JOIN offices AS parent_offices ON parent_offices.id = users.office_id AND parent_offices.level <= 3
	// 	WHERE positions.name = ? AND users.id <> ? AND users.deleted_at IS NULL
	// 	AND (parent_offices.id IS NOT NULL OR offices.level = 3)
	// `).Scan(&tos, typeID, typeValue, user.OfficeID, user.Office.ParentID, positionName, userID)

	return submission.RequirementDB{}, nil
}
