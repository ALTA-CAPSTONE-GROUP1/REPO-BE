package repository

import (
	"errors"
	"fmt"
	"strconv"

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

// GenerateIDFromPositionTag generates a unique user ID based on the position tag
func (um *usersModel) GenerateIDFromPositionTag(positionTag string) (string, error) {
	id := positionTag + "1"

	exists, err := um.CheckUserIDExists(id)
	if err != nil {
		return "", err
	}

	for exists {
		numberPart := id[len(positionTag):]

		number, err := strconv.Atoi(numberPart)
		if err != nil {
			return "", err
		}
		number++

		id = positionTag + strconv.Itoa(number)

		exists, err = um.CheckUserIDExists(id)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

// checkUserIDExists checks if a user ID already exists in the 'users' table
func (um *usersModel) CheckUserIDExists(id string) (bool, error) {
	var count int64
	err := um.db.Table("users").Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetPositionTagByID implements user.Repository
func (um *usersModel) GetPositionTagByID(positionID int) (string, error) {
	var position admin.Position
	err := um.db.Table("positions").Select("tag").Where("id = ?", positionID).First(&position).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("position tag not found for position ID: %d", positionID)
		}
		return "", err
	}

	return position.Tag, nil
}

// DeleteUser implements user.Repository
func (um *usersModel) DeleteUser(id string) error {
	tx := um.db.Table("Users").Where("id = ?", id).Delete(&admin.Users{})
	if tx.RowsAffected < 1 {
		log.Error("Terjadi error")
		return errors.New("no data deleted")
	}
	if tx.Error != nil {
		log.Error("User tidak ditemukan")
		return tx.Error
	}
	return nil
}

// UpdateUser implements user.Repository
func (um *usersModel) UpdateUser(id string, input user.Core) error {
	var updateUser admin.Users

	positionTag, err := um.GetPositionTagByID(input.PositionID)
	if err != nil {
		log.Error("error getting position tag", err.Error())
		return err
	}

	updatedID, err := um.GenerateIDFromPositionTag(positionTag)
	if err != nil {
		log.Error("error generating user ID", err.Error())
		return err
	}

	// Update the ID and other fields
	updateUser.ID = updatedID
	updateUser.Name = input.Name
	updateUser.Email = input.Email
	updateUser.PhoneNumber = input.PhoneNumber
	updateUser.OfficeID = input.OfficeID
	updateUser.PositionID = input.PositionID

	// Update the user in the database
	tx := um.db.Model(&admin.Users{}).Where("id = ?", id).Updates(&updateUser)
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
func (ur *usersModel) GetUserById(id string) (user.Core, error) {
	var user user.Core
	err := ur.db.Table("users").
		Joins("JOIN positions ON positions.id = users.position_id").
		Joins("JOIN offices ON offices.id = users.office_id").
		Select("users.id, users.office_id, users.position_id, users.name, users.email, users.phone_number, users.password, positions.name as position_name, offices.name as office_name").
		Where("users.id = ?", id).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		log.Error("failed to find user:", err.Error())
		return user, errors.New("failed to retrieve user")
	}

	return user, nil
}

// SelectAllUser implements user.Repository
func (ur *usersModel) SelectAllUser(limit, offset int, name string) ([]user.Core, error) {
	var users []user.Core
	var dbuser []admin.Users
	query := ur.db.Table("users").
		Joins("JOIN positions ON positions.id = users.position_id").
		Joins("JOIN offices ON offices.id = users.office_id").
		Select("users.id, users.office_id, users.position_id, users.name, users.email, users.phone_number, users.password, positions.name as position_name, offices.name as office_name").
		Limit(limit).
		Offset(offset)

	if name != "" {
		query = query.Where("users.name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&dbuser).Error
	if err != nil {
		log.Error("failed to find all users:", err.Error())
		return nil, errors.New("failed to retrieve users")
	}

	for _, v := range dbuser {
		tmp := user.Core{
			ID:          v.ID,
			OfficeID:    v.OfficeID,
			PositionID:  v.PositionID,
			Name:        v.Name,
			Email:       v.Email,
			PhoneNumber: v.Email,
		}
		users = append(users, tmp)
	}

	return users, nil
}

// InsertUser implements user.Repository
func (um *usersModel) InsertUser(newUser user.Core) error {
	inputUser := admin.Users{}

	positionTag, err := um.GetPositionTagByID(newUser.PositionID)
	if err != nil {
		log.Error("error getting position tag", err.Error())
		return err
	}

	id, err := um.GenerateIDFromPositionTag(positionTag)
	if err != nil {
		log.Error("error generating user ID", err.Error())
		return err
	}
	inputUser.ID = id

	samePassword := "kadal123"

	hashedPassword, err := helper.HashPassword(samePassword)
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

	existingUser := admin.Users{}
	if err := um.db.Where("email = ?", newUser.Email).First(&existingUser).Error; err == nil {
		return errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("error occurred while checking duplicate email", err.Error())
		return err
	}

	if err := um.db.Where("phone_number = ?", newUser.PhoneNumber).First(&existingUser).Error; err == nil {
		return errors.New("phone number already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("error occurred while checking duplicate phone number", err.Error())
		return err
	}

	if err := um.db.Create(&inputUser).Error; err != nil {
		log.Error("error on create table users", err.Error())
		return err
	}

	return nil
}
