package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
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

// SelectSubmissionAprrove implements approve.Repository
func (ar *approverModel) SelectSubmissionAprrove(limit, offset int, search string) ([]approve.Core, error) {
	var submissions []approve.Core

	query := ar.db.Table("Submissions").
		Joins("JOIN user_to ON Submissions.id = user_to.submission_id").
		Where("user_to.user_id = user_id").
		Limit(limit).
		Offset(offset)

	err := query.Find(&submissions).Error
	if err != nil {
		log.Error("failed to find all submission for actions", err.Error())
		return nil, errors.New("failed to retrieve submission")
	}

	return submissions, nil
}
