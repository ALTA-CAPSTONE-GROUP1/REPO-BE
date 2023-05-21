package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
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
	var dbsub user.Submission

	query := ar.db.
		Table("submissions").
		Select("submissions.title, users.name as user_name, positions.name as user_position").
		Joins("JOIN tos ON submissions.id = tos.submission_id").
		Joins("JOIN users ON tos.user_id = users.id").
		Joins("JOIN types ON submissions.type_id = types.id").
		Joins("JOIN positions ON positions.id = users.position_id").
		Where("submissions.id = ?", id).
		Preload("Type").
		Preload("User").
		Preload("Tos", "submission_id = ?", id).
		Preload("Ccs", "submission_id = ?", id).
		Preload("Signs", "submission_id = ?", id).
		Find(&dbsub)

	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return approve.Core{}, errors.New("submission not found")
		}
		log.Error("failed to find submission:", query.Error.Error())
		return approve.Core{}, errors.New("failed to retrieve submission")
	}

	var submissionTypeName string
	if token == config.TokenSuperAdmin {
		submissionTypeName = dbsub.Type.Name
	} else {
		return approve.Core{}, errors.New("invalid token")
	}

	coreData := approve.Core{
		ID:        dbsub.ID,
		UserID:    dbsub.UserID,
		TypeID:    dbsub.TypeID,
		Title:     dbsub.Title,
		Message:   dbsub.Message,
		Status:    dbsub.Status,
		Is_Opened: false,
		CreatedAt: dbsub.CreatedAt,
		Type:      subtype.Core{SubmissionTypeName: submissionTypeName},
	}

	for _, v := range dbsub.Tos {
		cTos := approve.ToCore{
			SubmissionID: v.SubmissionID,
			UserID:       v.User.ID,
			Action_Type:  v.Action_Type,
			User:         v.User.Name,
			Position:     v.User.Position.Name,
		}
		coreData.Tos = append(coreData.Tos, cTos)
	}
	for _, v := range dbsub.Files {
		cFiles := approve.FileCore{
			SubmissionID: v.SubmissionID,
			Name:         v.Name,
			Link:         v.Link,
		}
		coreData.Files = append(coreData.Files, cFiles)
	}

	return coreData, nil
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
