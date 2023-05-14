package usecase

import (
	"errors"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth"
	"github.com/labstack/gommon/log"
)

type authLogic struct {
	u auth.Repository
}

func New(r auth.Repository) auth.UseCase {
	return &authLogic{
		u: r,
	}
}

// LogInLogic implements auth.UseCase
func (ul *authLogic) LogInLogic(id string, password string) (auth.Core, error) {
	res, err := ul.u.Login(id, password)
	if err != nil {

		if strings.Contains(err.Error(), "not exist") {
			return auth.Core{}, errors.New("id cannot be blank")

		} else if strings.Contains(err.Error(), "wrong") {
			return auth.Core{}, errors.New("password is wrong")

		}
		log.Error("error on loginlogic, internal server error", err.Error())
		return auth.Core{}, errors.New("internal server error")

	}

	return res, nil
}
