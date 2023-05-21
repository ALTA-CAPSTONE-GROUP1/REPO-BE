package usecase

import (
	"errors"
	"strings"

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

// UpdateUser implements approve.UseCase
func (al *approverLogic) UpdateApprove(userID string, id int, updateInput approve.Core) error {
	if err := al.a.UpdateApprove(userID, id, updateInput); err != nil {
		log.Error("failed on calling updateprofile query")
		if strings.Contains(err.Error(), "hashing password") {
			log.Error("hashing password error")
			return errors.New("is invalid")
		} else if strings.Contains(err.Error(), "affected") {
			log.Error("no rows affected on update submission")
			return errors.New("data is up to date")
		}
		return err
	}
	return nil
}

// GetSubmissionById implements approve.UseCase
func (al *approverLogic) GetSubmissionById(userID string, id int) (approve.Core, error) {
	result, err := al.a.SelectSubmissionById(userID, id)
	if err != nil {
		log.Error("failed to find submission for action", err.Error())
		return approve.Core{}, errors.New("internal server error")
	}

	return result, nil
}

// GetSubmissionAprrove implements approve.UseCase
func (al *approverLogic) GetSubmissionAprrove(userID string, search approve.GetAllQueryParams) ([]approve.Core, error) {
	result, err := al.a.SelectSubmissionAprrove(userID, search)
	if err != nil {
		log.Error("failed to find submission for action", err.Error())
		return []approve.Core{}, errors.New("internal server error")
	}

	return result, nil
}
