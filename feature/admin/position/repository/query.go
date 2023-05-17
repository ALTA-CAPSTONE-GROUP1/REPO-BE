package repository

import (
	"fmt"

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

func (pm *positionModel) InsertPosition(position position.Core) error {
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

func (pm *positionModel) GetPositions(limit int, offset int, search string) ([]position.Core, int64, error) {
	var (
		positions   []position.Core
		DBpositions []admin.Position
		count       int64
	)

	tx := pm.db.Find(&DBpositions)
	if tx.Error != nil {
		log.Error("get posititons query error without search condition")
		return nil, 0, tx.Error
	}
	tx.Count(&count)

	for _, dbPos := range DBpositions {
		corePos := position.Core{
			ID:   dbPos.ID,
			Name: dbPos.Name,
			Tag:  dbPos.Tag,
		}
		positions = append(positions, corePos)
	}

	return positions, count, nil
}

func (pm *positionModel) DeletePosition(position int) error {

	tx := pm.db.Where("id = ?", position).Delete(&admin.Position{})
	if tx.Error != nil {
		log.Error("delete position query error")
		return fmt.Errorf("delete position query error: %w", tx.Error)
	}
	return nil
}
