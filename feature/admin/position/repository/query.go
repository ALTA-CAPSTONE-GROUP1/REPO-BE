package repository

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type positionModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) position.Repository {
	return &positionModel{
		db: db,
	}
}

func (pm positionModel) InsertPositionHandler(position position.Core) error {
	newPosition := admin.Position{
		Name: position.Name,
		Tag:  position.Tag,
	}

	tx := pm.db.Create(&newPosition)
	if tx.Error != nil {
		log.Errorf("error on insert data %s with tag %s to position table in db", position.Name, position.Tag)
		return tx.Error
	}

	return nil
}
