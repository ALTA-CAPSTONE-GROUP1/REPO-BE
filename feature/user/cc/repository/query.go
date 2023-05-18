package repository

import (
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
	var submissions []subRepo.Cc

	if err := cm.db.Where("user_id = ?", userID).Find(&ccsOwned).Error; err != nil {
		log.Errorf("error on finding cc by userid %w", err)
		return []cc.CcCore{}, err
	}

	if err := cm.db.Where("user_id = ?", userID).
		Preload("Files").
		Preload("Tos").
		Preload("Ccs").
		Preload("Signs").
		Find(&submissions).
		Error; err != nil {
		log.Errorf("error on finding submissions for user %s: %v", userID, err)
		return []cc.CcCore{}, err
	}
	
}
