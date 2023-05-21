package usecase

import (
	"errors"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
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

// GetSubmissionById implements approve.UseCase
func (al *approverLogic) GetSubmissionByHyperApproval(userID string, id int, token string) (approve.GetSubmissionByIDCore, error) {
	result, err := al.a.SelectSubmissionByHyperApproval(userID, id, token)
	if err != nil {
		log.Error("failed to find submission for action", err.Error())
		return approve.GetSubmissionByIDCore{}, errors.New("internal server error")
	}

	return result, nil

}

// UpdateUser implements approve.UseCase
func (al *approverLogic) UpdateByHyperApproval(userID string, updateInput approve.Core) error {
	if err := al.a.UpdateByHyperApproval(userID, updateInput); err != nil {
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
