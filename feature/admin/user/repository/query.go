package repository

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type usersModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.Repository {
	return &usersModel{
		db: db,
	}
}

// InsertUser implements user.Repository
func (um *usersModel) InsertUser(newUser user.Core) error {
	inputUser := admin.Users{}
	hashedPassword, err := helper.HashPassword(newUser.Password)
	if err != nil {
		log.Error("error occurs on hashing password", err.Error())
		return err
	}

	inputUser.Name = newUser.Name
	inputUser.Email = newUser.Email
	inputUser.PhoneNumber = newUser.PhoneNumber
	inputUser.OfficeID = newUser.OfficeID
	inputUser.PositionID = newUser.PositionID
	inputUser.Password = hashedPassword

	if err := um.db.Create(&inputUser).Error; err != nil {
		log.Error("error on create table users", err.Error())
		return err
	}

	return nil
}
