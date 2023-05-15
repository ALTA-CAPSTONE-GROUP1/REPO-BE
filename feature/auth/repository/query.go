package repository

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type authModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) auth.Repository {
	return &authModel{
		db: db,
	}
}

// Login implements auth.Repository
func (um *authModel) Login(id string, password string) (auth.Core, error) {
	input := User{}

	if id == "" {
		log.Error("id login is blank")
		return auth.Core{}, errors.New("data does not exist")
	}

	if err := um.db.Where("id = ?", id).First(&input).Error; err != nil {
		log.Error("error occurs on select users login", err.Error())
		return auth.Core{}, err
	}

	// Check if the user is the admin and the password is "admin"
	if input.ID == "admin" && password == "admin" {
		return auth.Core{
			ID:          input.ID,
			OfficeID:    input.OfficeID,
			PositionID:  input.PositionID,
			Name:        input.Name,
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
			Password:    input.Password,
		}, nil
	}

	// Perform the regular password verification for non-admin users
	if err := helper.VerifyPassword(input.Password, password); err != nil {
		log.Error("user input for password is wrong", err.Error())
		return auth.Core{}, errors.New("wrong password")
	}

	return auth.Core{
		ID:          input.ID,
		OfficeID:    input.OfficeID,
		PositionID:  input.PositionID,
		Name:        input.Name,
		Email:       input.Email,
		PhoneNumber: input.PhoneNumber,
		Password:    input.Password,
	}, nil
}
