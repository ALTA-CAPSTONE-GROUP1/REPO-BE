package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type userModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) profile.Repository {
	return &userModel{
		db: db,
	}
}

// Profile implements profile.Repository
func (um *userModel) Profile(id string) (profile.Core, error) {
	tmp := profile.Core{}
	tx := um.db.Table("users").Where("id = ?", id).First(&tmp)
	if tx.RowsAffected < 1 {
		log.Error("Terjadi error saat first user (data tidak ditemukan)")
		return profile.Core{}, errors.New("user not found")
	}
	if tx.Error != nil {
		log.Error("Terjadi Kesalahan")
		return profile.Core{}, tx.Error
	}
	return tmp, nil
}
