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
