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
		positions []position.Core
		count     int64
	)
	qry := pm.db.Limit(limit).Offset(offset).Select("id, name, tag").Table("positions").Order("id DESC")

	if search != "" {
		search = "%" + search + "%"
		if err := qry.Where("name LIKE ? OR tag LIKE ?", search, search).Find(&positions).Error; err != nil {
			log.Errorf("error on finding users with search param")
			return nil, 0, err
		}
		if err := pm.db.Where("name LIKE ? OR tag LIKE ?", search, search).Select("id, name").Table("positions").Count(&count).Error; err != nil {
			log.Errorf("counting total data %w", err)
			return nil, 0, err
		}
	} else {
		if err := qry.Find(&positions).Error; err != nil {
			log.Errorf("error on finding user without search params %w", err)
			return nil, 0, err
		}
		if err := pm.db.Select("id, name, tag").Table("positions").Count(&count).Error; err != nil {
			log.Errorf("counting total data %w", err)
			return nil, 0, err
		}
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
