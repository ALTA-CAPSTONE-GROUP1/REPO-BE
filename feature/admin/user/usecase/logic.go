package usecase

import (
	"errors"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
	"github.com/labstack/gommon/log"
)

type userLogic struct {
	u user.Repository
}

func New(u user.Repository) user.UseCase {
	return &userLogic{
		u: u,
	}
}

// RegisterUser implements user.UseCase
func (ul *userLogic) RegisterUser(newUser user.Core) error {
	if err := ul.u.InsertUser(newUser); err != nil {
		log.Error("error on calling register insert user query", err.Error())
		if strings.Contains(err.Error(), "column") {
			return errors.New("server error")
		} else if strings.Contains(err.Error(), "value") {
			return errors.New("invalid value")
		} else if strings.Contains(err.Error(), "too short") {
			return errors.New("invalid password length")
		}
		return errors.New("server error")
	}
	return nil
}
