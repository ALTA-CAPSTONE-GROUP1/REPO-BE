package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
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
	var dbuser user.Users
	err := um.db.Table("users").
		Preload("Position").
		Preload("Office").
		Joins("JOIN positions ON positions.id = users.position_id").
		Joins("JOIN offices ON offices.id = users.office_id").
		Select("users.id, users.office_id, users.position_id, users.name, users.email, users.phone_number, users.password, positions.name as position_name, offices.name as office_name").
		Where("users.id = ?", id).
		First(&dbuser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return profile.Core{}, errors.New("user not found")
		}
		log.Error("failed to find user:", err.Error())
		return profile.Core{}, errors.New("failed to retrieve user")
	}

	return profile.Core{
		ID:          dbuser.ID,
		OfficeID:    dbuser.OfficeID,
		PositionID:  dbuser.PositionID,
		Name:        dbuser.Name,
		Email:       dbuser.Email,
		PhoneNumber: dbuser.PhoneNumber,
		Password:    dbuser.Password,
		Position:    position.Core{Name: dbuser.Position.Name},
		Office:      office.Core{Name: dbuser.Office.Name},
	}, nil
}
