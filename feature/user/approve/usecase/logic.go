package usecase

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	"github.com/labstack/gommon/log"
)

type approverLogic struct {
	a approve.Repository
}

func New(a approve.Repository) approve.UseCase {
	return &approverLogic{
		a: a,
	}
}

// GetSubmissionAprrove implements approve.UseCase
func (al *approverLogic) GetSubmissionAprrove(limit, offset int, search string) ([]approve.Core, error) {
	result, err := al.a.SelectSubmissionAprrove(limit, offset, search)
	if err != nil {
		log.Error("failed to find submission for action", err.Error())
		return []approve.Core{}, errors.New("internal server error")
	}

	return result, nil
}
