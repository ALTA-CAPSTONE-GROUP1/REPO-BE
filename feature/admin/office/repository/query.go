package repository

import (
	"errors"

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

// DeleteOffice implements office.Repository
func (om *officeModel) DeleteOffice(id uint) error {
	tx := om.db.Where("id = ?", id).Delete(&admin.Office{})
	if tx.RowsAffected < 1 {
		log.Error("Terjadi error")
		return errors.New("no data deleted")
	}
	if tx.Error != nil {
		log.Error("Office tidak ditemukan")
		return tx.Error
	}
	return nil
}

// GetAllOffice implements office.Repository
func (om *officeModel) GetAllOffice(limit int, offset int, search string) ([]office.Core, int, error) {
	nameSearch := "%" + search + "%"
	totalData := int64(-1)
	var res []office.Core

	qry := om.db.Limit(limit).Offset(offset).Select("id, name").Table("offices").Order("id DESC")

	if search != "" {
		if err := qry.Where("name LIKE ?", nameSearch).Find(&res).Error; err != nil {
			log.Errorf("error on finding search %w", err)
			return []office.Core{}, int(totalData), nil
		}
		if err := om.db.Where("offices.name LIKE = ?", nameSearch).Select("id, name").Table("offices").Count(&totalData).Error; err != nil {
			log.Errorf("error on count filtered data %w", err)
			return []office.Core{}, int(totalData), nil
		}
	} else {
		if err := qry.Find(&res).Error; err != nil {
			log.Errorf("error on finding data without search %w", err)
			return []office.Core{}, int(totalData), nil
		}
		if err := om.db.Select("id, name").Table("offices").Count(&totalData).Error; err != nil {
			log.Errorf("error on counting data without search %w", err)
			return []office.Core{}, int(totalData), nil
		}
	}

	return res, int(totalData), nil
}

// InsertOffice implements office.Repository
func (om *officeModel) InsertOffice(newOffice office.Core) error {
	inputOffice := admin.Office{
		Name: newOffice.Name,
	}

	tx := om.db.Create(&inputOffice)
	if tx.Error != nil {
		log.Error("error on create table office")
		return tx.Error
	}

	return nil
}
