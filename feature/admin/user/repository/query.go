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
	id := positionTag + "01"

	// Check if the generated ID already exists in the 'users' table
	exists, err := um.CheckUserIDExists(id)
	if err != nil {
		return "", err
	}

	// If the ID exists, increment the number part and check again until a unique ID is found
	for exists {
		// Extract the number part from the ID
		numberPart := id[len(positionTag):]

		// Parse the number part and increment it
		number, err := strconv.Atoi(numberPart)
		if err != nil {
			return "", err
		}
		number++

		// Generate the updated ID by combining the position tag and the incremented number
		id = positionTag + strconv.Itoa(number)

		// Check if the updated ID exists
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
	query := ur.db.Table("users").
		Joins("JOIN positions ON positions.id = users.position_id").
		Joins("JOIN offices ON offices.id = users.office_id").
		Select("users.id, users.office_id, users.position_id, users.name, users.email, users.phone_number, users.password, positions.name as position_name, offices.name as office_name").
		Limit(limit).
		Offset(offset)

	if name != "" {
		query = query.Where("users.name LIKE ?", "%"+name+"%")
	}

	err := query.Find(&users).Error
	if err != nil {
		log.Error("failed to find all users:", err.Error())
		return nil, errors.New("failed to retrieve users")
	}

	return users, nil
}

// InsertUser implements user.Repository
func (um *usersModel) InsertUser(newUser user.Core) error {
	inputUser := admin.Users{}

	// Get the position tag based on newUser.PositionID
	positionTag, err := um.GetPositionTagByID(newUser.PositionID)
	if err != nil {
		log.Error("error getting position tag", err.Error())
		return err
	}

	// Generate the inputUser.ID based on the position tag
	id, err := um.GenerateIDFromPositionTag(positionTag)
	if err != nil {
		log.Error("error generating user ID", err.Error())
		return err
	}
	inputUser.ID = id

	// Rest of the code to hash the password and populate the inputUser fields
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
