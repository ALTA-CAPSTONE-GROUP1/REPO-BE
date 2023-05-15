package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
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

// UpdateUser implements profile.Repository
func (um *userModel) UpdateUser(id string, input profile.Core) error {
	var updateUser user.Users

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		log.Error("error occurs on hashing password", err.Error())
		return errors.New("hashing password failed")
	}

	updateUser.Email = input.Email
	updateUser.PhoneNumber = input.PhoneNumber
	updateUser.Password = hashedPassword

	// Update the user in the database
	tx := um.db.Model(&user.Users{}).Where("id = ?", id).Updates(&updateUser)
	if tx.RowsAffected < 1 {
		log.Error("there is no column to change on update user")
		return errors.New("no data affected")
	}
	if tx.Error != nil {
		log.Error("error on update user")
		return tx.Error
	}

	return nil
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
