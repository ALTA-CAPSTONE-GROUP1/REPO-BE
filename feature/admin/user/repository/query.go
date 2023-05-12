package repository

import (
	"errors"

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

// UpdateUser implements user.Repository
func (um *usersModel) UpdateUser(id string, input user.Core) error {
	var UpdateUser admin.Users

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		log.Error("error occurs on hashing password", err.Error())
		return errors.New("hashing password failed")
	}

	UpdateUser.ID = input.ID
	UpdateUser.Name = input.Name
	UpdateUser.Email = input.Email
	UpdateUser.PhoneNumber = input.PhoneNumber
	UpdateUser.Password = hashedPassword
	UpdateUser.OfficeID = input.OfficeID
	UpdateUser.PositionID = input.PositionID

	tx := um.db.Model(&admin.Users{}).Where("id = ?", id).Updates(&UpdateUser)
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

// GetUserById implements user.Repository
func (um *usersModel) GetUserById(id string) (user.Core, error) {
	var res user.Core
	if err := um.db.Where("id = ?", id).First(&res).Error; err != nil {
		log.Error("error occurs in finding user profile", err.Error())
		return user.Core{}, err
	}

	return res, nil
}

// SelectAllUser implements user.Repository
func (um *usersModel) SelectAllUser(limit int, offset int, name string) ([]user.Core, error) {
	nameSearch := "%" + name + "%"
	var res []user.Core
	if err := um.db.Limit(limit).Offset(offset).Where("users.name LIKE ?", nameSearch).Select("users.id, users.email, users.phone_number, users.office_id, users.position_id").Find(&res).Error; err != nil {
		log.Error("error occurs in finding all user", err.Error())
		return nil, err
	}

	return res, nil
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
