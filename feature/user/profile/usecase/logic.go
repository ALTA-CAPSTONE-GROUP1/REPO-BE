package usecase

import (
	"errors"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile"
	"github.com/labstack/gommon/log"
)

type userLogic struct {
	u profile.Repository
}

func New(u profile.Repository) profile.UseCase {
	return &userLogic{
		u: u,
	}
}

// UpdateUser implements profile.UseCase
func (ul *userLogic) UpdateUser(id string, updateUser profile.Core) error {
	if err := ul.u.UpdateUser(id, updateUser); err != nil {
		log.Error("failed on calling updateprofile query")
		if strings.Contains(err.Error(), "hashing password") {
			log.Error("hashing password error")
			return errors.New("is invalid")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on update user")
			return errors.New("data is up to date")
		}
		return err
	}
	return nil
}

// ProfileLogic implements profile.UseCase
func (ul *userLogic) ProfileLogic(id string) (profile.Core, error) {
	result, err := ul.u.Profile(id)
	if err != nil {
		log.Error("failed to find user", err.Error())
		return profile.Core{}, errors.New("internal server error")
	}

	return result, nil
}
