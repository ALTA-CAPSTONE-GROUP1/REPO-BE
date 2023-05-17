package repository

import (
	"errors"
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
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

// SelectSubmissionById implements approve.Repository
func (ar *approverModel) SelectSubmissionById(userID string, id int) (approve.Core, error) {
	var res approve.Core
	var dbsub user.Submission

	query := ar.db.Table("submissions").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN types ON submissions.type_id = types.id").
		Where("tos.user_id = ? AND submissions.id = ?", userID, id).
		Preload("Type").
		Find(&dbsub)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return approve.Core{}, errors.New("submission not found")
		}
		log.Error("failed to find submission:", query.Error.Error())
		return approve.Core{}, errors.New("failed to retrieve submission")
	}

	res = approve.Core{
		ID:        dbsub.ID,
		UserID:    dbsub.UserID,
		TypeID:    dbsub.TypeID,
		Title:     dbsub.Title,
		Message:   dbsub.Message,
		Status:    dbsub.Status,
		Is_Opened: false,
		CreatedAt: time.Time{},
		Type: admin.Type{
			Name: dbsub.Type.Name,
		},
	}

	return res, nil
}

// SelectSubmissionApprove implements approve.Repository
func (ar *approverModel) SelectSubmissionAprrove(userID string, limit, offset int, search string) ([]approve.Core, error) {
	var res []approve.Core
	var dbsub []user.Submission

	query := ar.db.Table("submissions").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN types ON submissions.type_id = types.id").
		Where("tos.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Preload("Type").
		Find(&dbsub)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return []approve.Core{}, errors.New("submission not found")
		}
		log.Error("failed to find all submission:", query.Error.Error())
		return []approve.Core{}, errors.New("failed to retrieve all submission")
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
