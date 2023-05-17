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
func (om *officeModel) GetAllOffice(limit int, offset int, search string) ([]office.Core, error) {
	nameSearch := "%" + search + "%"
	var dbres []admin.Office
	var res []office.Core
	if err := om.db.Limit(limit).Offset(offset).Where("offices.name LIKE ?", nameSearch).Select("offices.id, offices.name").Find(&dbres).Error; err != nil {
		log.Error("error occurs in finding all office", err.Error())
		return nil, err
	}

	for _, v := range dbres {
		tmp := office.Core{
			ID:   v.ID,
			Name: v.Name,
		}
		res = append(res, tmp)
	}

	return res, nil
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
