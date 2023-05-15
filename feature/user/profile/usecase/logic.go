package usecase

import (
	"errors"

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

// ProfileLogic implements profile.UseCase
func (ul *userLogic) ProfileLogic(id string) (profile.Core, error) {
	result, err := ul.u.Profile(id)
	if err != nil {
		log.Error("failed to find user", err.Error())
		return profile.Core{}, errors.New("internal server error")
	}

	return result, nil
}
