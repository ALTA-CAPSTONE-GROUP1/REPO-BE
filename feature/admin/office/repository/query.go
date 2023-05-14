package repository

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type officeModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) office.Repository {
	return &officeModel{
		db: db,
	}
}

// GetAllOffice implements office.Repository
func (om *officeModel) GetAllOffice(limit int, offset int, search string) ([]office.Core, error) {
	nameSearch := "%" + search + "%"
	var res []office.Core
	if err := om.db.Limit(limit).Offset(offset).Where("offices.name LIKE ?", nameSearch).Select("offices.id, offices.name").Find(&res).Error; err != nil {
		log.Error("error occurs in finding all office", err.Error())
		return nil, err
	}

	return res, nil
}

// InsertOffice implements office.Repository
func (om *officeModel) InsertOffice(newOffice office.Core) error {
	inputOffice := admin.Office{}

	inputOffice.Name = newOffice.Name

	tx := om.db.Create(&newOffice)
	if tx.Error != nil {
		log.Error("error on create table office")
		return tx.Error
	}

	return nil
}
