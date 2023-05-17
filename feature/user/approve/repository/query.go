package repository

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"

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
// func (ar *approverModel) SelectSubmissionAprrove(userID string, limit, offset int, search string) ([]approve.Core, error) {
// 	var submissions []approve.Core

// 	query := ar.db.Table("submissions").
// 		Joins("JOIN tos ON submissions.id = tos.submission_id").
// 		Joins("JOIN users ON tos.user_id = users.id").
// 		Where("users.position_id = ? OR users.id IN (SELECT user_id FROM tos)", userID).
// 		Limit(limit).
// 		Offset(offset).
// 		Find(&submissions)
// 	if query.Error != nil {
// 		return nil, query.Error
// 	}

// 	return submissions, nil
// }

func (ar *approverModel) SelectSubmissionAprrove(userID string, limit, offset int, search string) ([]approve.Core, error) {
	var res []approve.Core
	var dbsub []user.Submission

	query := ar.db.Table("submissions").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Where("tos.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&res)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, v := range dbsub {
		tmp := approve.Core{
			ID: v.ID,
			// UserID:    v.userID,
			TypeID:    v.TypeID,
			Title:     v.Title,
			Message:   v.Message,
			Status:    v.Status,
			Is_Opened: false,
			Type:      admin.Type{},
			Files:     []user.File{},
			Tos:       []user.To{},
			Ccs:       []user.Cc{},
			Signs:     []user.Sign{},
		}
		res = append(res, tmp)
	}

	return res, nil

}
