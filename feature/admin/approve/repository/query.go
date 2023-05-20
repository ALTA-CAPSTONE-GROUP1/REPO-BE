package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
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

// / SelectSubmissionByHyperApproval implements approve.Repository
func (ar *approverModel) SelectSubmissionByHyperApproval(userID string, id int, token string) (approve.Core, error) {
	panic("unimplemented")
}

// UpdateUser implements approve.Repository
func (ar *approverModel) UpdateByHyperApproval(userID string, input approve.Core) error {
	// submission := approve.Core{}
	dbsub := user.Submission{}
	tx := ar.db.Model(&user.Submission{}).
		Where("id = ?", input.ID).
		First(&dbsub)

	if tx.RowsAffected == 0 {
		log.Error("no rows found for the given submission ID")
		return errors.New("no data found")
	}

	switch input.Status {
	case "approve":
		input.Status = "waiting"
	case "revise":
		input.Status = "revised"
	case "reject":
		input.Status = "rejected"
	default:
		return errors.New("invalid status")
	}

	tx = ar.db.Model(&dbsub).
		Where("id = ?", input.ID).
		Updates(user.Submission{Status: input.Status})

	if tx.RowsAffected == 0 {
		log.Error("no rows affected on update submission")
		return errors.New("data is up to date")
	}

	if tx.Error != nil {
		log.Error("error on update submission")
		return tx.Error
	}

	to := user.To{}
	tx = ar.db.Model(&user.To{}).
		Where("submission_id = ?", input.ID).
		First(&to)

	if tx.RowsAffected == 0 {
		log.Error("no rows found for the given submission ID in user.To")
		return errors.New("no data found in user.To")
	}

	tx = ar.db.Model(&user.To{}).
		Where("submission_id = ?", input.ID).
		Updates(approve.ToCore{Message: config.AutoMessageHyperApp})

	if tx.RowsAffected == 0 {
		log.Error("no rows affected on update to message")
		return errors.New("data is up to date")
	}

	if tx.Error != nil {
		log.Error("error on update to message")
		return tx.Error
	}

	actionType := ""
	switch input.Status {
	case "waiting":
		actionType = "approve"
	case "revised":
		actionType = "revise"
	case "rejected":
		actionType = "reject"
	}

	tx = ar.db.Model(&user.To{}).
		Joins("JOIN users ON user_id = users.id").
		Where("submission_id = ?", input.ID).
		Update("action_type", actionType)

	if tx.Error != nil {
		log.Error("error on update action_type in 'to' table")
		return tx.Error
	}

	return nil
}
