package repository

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
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
	tx := um.db.Where("id = ?", id).Delete(&admin.Users{})
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

	// Fetch the existing user from the database
	existingUser := admin.Users{}
	err := um.db.Where("id = ?", id).First(&existingUser).Error
	if err != nil {
		log.Error("error fetching existing user", err.Error())
		return err
	}

	// Update the ID and other fields
	updateUser.ID = existingUser.ID
	updateUser.Name = input.Name
	updateUser.Email = input.Email
	updateUser.PhoneNumber = input.PhoneNumber

	// Check if the position ID is valid
	if input.PositionID != 0 {
		positionTag, err := um.GetPositionTagByID(input.PositionID)
		if err != nil {
			log.Error("error getting position tag", err.Error())
			return err
		}
		updateUser.PositionID = input.PositionID
		updatedID, err := um.GenerateIDFromPositionTag(positionTag)
		if err != nil {
			log.Error("error generating user ID", err.Error())
			return err
		}
		updateUser.ID = updatedID
	} else {
		updateUser.PositionID = existingUser.PositionID
	}

	// Check if the office ID is valid
	if input.OfficeID != 0 {
		updateUser.OfficeID = input.OfficeID
	} else {
		updateUser.OfficeID = existingUser.OfficeID
	}

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
	// var user user.Core
	var dbuser admin.Users
	err := ur.db.Table("users").
		Preload("Position").
		Preload("Office").
		Joins("JOIN positions ON positions.id = users.position_id").
		Joins("JOIN offices ON offices.id = users.office_id").
		Select("users.id, users.office_id, users.position_id, users.name, users.email, users.phone_number, users.password, positions.name as position_name, offices.name as office_name").
		Where("users.id = ?", id).
		First(&dbuser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user.Core{}, errors.New("user not found")
		}
		log.Error("failed to find user:", err.Error())
		return user.Core{}, errors.New("failed to retrieve user")
	}

	return user.Core{
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

// SelectAllUser implements user.Repository
func (ur *usersModel) SelectAllUser(limit, offset int, name string) ([]user.Core, int, error) {
	var users []user.Core
	var dbusers []admin.Users
	nameSearch := "%" + name + "%"
	totalData := int64(-1)

	query := ur.db.Table("users").
		Preload("Position").
		Preload("Office").
		Joins("JOIN positions ON positions.id = users.position_id").
		Joins("JOIN offices ON offices.id = users.office_id").
		Select("users.id, users.office_id, users.position_id, users.name, users.email, users.phone_number, users.password, positions.name as position_name, offices.name as office_name").
		Limit(limit).
		Offset(offset).
		Order("id DESC")

	if name != "" {
		if err := query.Where("users.name LIKE ? OR users.email LIKE ? OR users.phone_number LIKE ? OR positions.name LIKE ? OR offices.name LIKE ?",
			nameSearch, nameSearch, nameSearch, nameSearch, nameSearch).Find(&dbusers).Error; err != nil {
			log.Errorf("error on finding search: %w", err)
			return []user.Core{}, int(totalData), err
		}
		if err := ur.db.Table("users").
			Joins("JOIN positions ON positions.id = users.position_id").
			Joins("JOIN offices ON offices.id = users.office_id").
			Where("users.name LIKE ? OR users.email LIKE ? OR users.phone_number LIKE ? OR positions.name LIKE ? OR offices.name LIKE ?",
				nameSearch, nameSearch, nameSearch, nameSearch, nameSearch).
			Count(&totalData).Error; err != nil {
			log.Errorf("error on count filtered data: %w", err)
			return []user.Core{}, int(totalData), err
		}
	} else {
		if err := query.Find(&dbusers).Error; err != nil {
			log.Errorf("error on finding data without search: %w", err)
			return []user.Core{}, int(totalData), err
		}
		if err := ur.db.Table("users").Count(&totalData).Error; err != nil {
			log.Errorf("error on counting data without search: %w", err)
			return []user.Core{}, int(totalData), err
		}
	}

	for _, dbuser := range dbusers {
		position := position.Core{
			Name: dbuser.Position.Name,
		}
		office := office.Core{
			Name: dbuser.Office.Name,
		}

		user := user.Core{
			ID:          dbuser.ID,
			OfficeID:    dbuser.OfficeID,
			PositionID:  dbuser.PositionID,
			Name:        dbuser.Name,
			Email:       dbuser.Email,
			PhoneNumber: dbuser.PhoneNumber,
			Password:    dbuser.Password,
			Position:    position,
			Office:      office,
		}

		users = append(users, user)
	}

	return users, int(totalData), nil
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

	samePassword := "alta123"

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

	// registrationData := fmt.Sprintf("Name: %s\nEmail: %s\nPhone Number: %s\nOffice: %s\nPosition: %s\nPassword: %s",
	// 	newUser.Name,
	// 	newUser.Email,
	// 	newUser.PhoneNumber,
	// 	newUser.Office.Name,
	// 	newUser.Position.Name,
	// 	samePassword,
	// )

	// if err := helper.SendTelegramMessage(newUser.PhoneNumber, registrationData); err != nil {
	// 	log.Error("failed to send notifications via Telegram:", err.Error())
	// }

	// if err := helper.SendSimpleEmail("Registration Update", "Registration Update", []string{newUser.Email}, []string{newUser.Name}, "eProposal Account"); err != nil {
	// 	log.Error("failed to send email:", err.Error())
	// }

	return nil

}

// if err := sm.db.Preload("Position").
//         Joins("INNER JOIN position_has_types ON position_has_types.position_id = users.position_id").
//         Joins("INNER JOIN positions ON positions.id = position_has_types.position_id").
//         Joins("INNER JOIN offices ON offices.id = users.office_id").
//         Joins("INNER JOIN types ON types.id = position_has_types.type_id").
//         Where("positions.deleted_at IS NULL AND types.deleted_at IS NULL AND position_has_types.deleted_at IS NULL").
//         Where("types.name = ? AND position_has_types.as = 'To' AND position_has_types.value = ? AND (users.office_id = ? OR offices.name = 'Head Office')",
// 		typeName, typeValue, applicant.Office.ID).
//         Order("position_has_types.to_level ASC").

// 		if err := sm.db.Preload("Position").
//         Joins("INNER JOIN position_has_types ON position_has_types.position_id = users.position_id").
//         Joins("INNER JOIN positions ON positions.id = position_has_types.position_id").
//         Joins("INNER JOIN offices ON offices.id = users.office_id").
//         Joins("INNER JOIN types ON types.id = position_has_types.type_id").
//         Where("positions.deleted_at IS NULL AND types.deleted_at IS NULL AND position_has_types.deleted_at IS NULL").
//         Where("types.name = ? AND position_has_types.as = 'To' AND position_has_types.value = ? AND
// 		(users.office_id = ? OR offices.name = 'Head Office')", typeName, typeValue, applicant.Office.ID).
//         Order("position_has_types.to_level ASC").
//         Find(&tos).Error; err != nil {
//         return submission.Core{}, err
//     }
