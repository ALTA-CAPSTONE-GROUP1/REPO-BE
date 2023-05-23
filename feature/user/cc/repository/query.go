package repository

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	subRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/repository"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type ccModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) cc.Repository {
	return &ccModel{
		db: db,
	}
}

func (cm *ccModel) GetAllCc(userID string) ([]cc.CcCore, error) {
	var result []cc.CcCore
	var ccsOwned []subRepo.Cc
	var submissions []subRepo.Submission

	if err := cm.db.Where("user_id = ?", userID).Find(&ccsOwned).Error; err != nil {
		log.Errorf("error on finding cc by userid %w", err)
		return []cc.CcCore{}, err
	}

	for _, ccOwned := range ccsOwned {
		if err := cm.db.Where("id = ?", ccOwned.SubmissionID).
			Preload("Files").
			Preload("Tos").
			Preload("Ccs").
			Find(&submissions).
			Error; err != nil {
			log.Errorf("error on finding submissions for user %s: %v", userID, err)
			return []cc.CcCore{}, err
		}

		for _, submission := range submissions {
			if submission.Status != "Approved" {
				continue
			}
			var toUser admin.Users
			if err := cm.db.Where("id = ?", submission.Tos[0].UserID).Preload("Position").First(&toUser).Error; err != nil {
				log.Errorf("error on finding to user %w", err)
				return []cc.CcCore{}, err
			}

			var fromUser admin.Users
			if err := cm.db.Where("id = ?", submission.UserID).Preload("Position").First(&fromUser).Error; err != nil {
				log.Errorf("error on finding to user %w", err)
				return []cc.CcCore{}, err
			}

			var subType admin.Type
			if err := cm.db.Where("id = ?", submission.TypeID).First(&subType).Error; err != nil {
				log.Errorf("error on finding submission type detail %w", err.Error())
				return []cc.CcCore{}, err
			}

			tmp := cc.CcCore{
				SubmisisonID: ccOwned.SubmissionID,
				Title:        submission.Title,
				Attachment:   submission.Files[0].Link,
				To: cc.Receiver{
					Name:     toUser.Name,
					Position: toUser.Position.Name,
				},
				From: cc.Sender{
					Name:     fromUser.Name,
					Position: fromUser.Position.Name,
				},
				SubmissionType: subType.Name,
			}

			result = append(result, tmp)
		}
	}
	return result, nil
}

// SELECT submissions.id, submissions.title, submissions.type_id, files.link, tos.id as to_id , `types`.`name`, users.`id`, positions.`name`
// 		FROM submissions
// 		JOIN files ON submissions.id = files.submission_id
// 		JOIN tos ON submissions.id = tos.submission_id
// 		JOIN ccs ON submissions.id = ccs.submission_id
//         JOIN `types` ON submissions.type_id = `types`.id
//         JOIN users ON submissions.user_id = users.id
//         JOIN positions ON users.position_id = positions.id
// 		WHERE submissions.id IN (7,8) AND ccs.user_id = 'AASpv1'
// 		ORDER BY submissions.id ASC
// 		LIMIT 10
// 		OFFSET 0;
