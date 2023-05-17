package repository

import (
	"time"

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

func (ar *approverModel) SelectSubmissionAprrove(userID string, limit, offset int, search string) ([]approve.Core, error) {
	var res []approve.Core
	var dbsub []user.Submission

	query := ar.db.Table("submissions").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Where("tos.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&dbsub)
	if query.Error != nil {
		return nil, query.Error
	}
	for _, v := range dbsub {
		tmp := approve.Core{
			ID:        v.ID,
			UserID:    v.UserID,
			TypeID:    v.TypeID,
			Title:     v.Title,
			Message:   v.Message,
			Status:    v.Status,
			Is_Opened: false,
			CreatedAt: time.Time{},
			Type: admin.Type{
				Name: v.Type.Name,
			},
		}
		res = append(res, tmp)
	}

	return res, nil

}
