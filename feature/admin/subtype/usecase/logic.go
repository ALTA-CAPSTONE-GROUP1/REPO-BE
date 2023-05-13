package usecase

import (
	"errors"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/labstack/gommon/log"
)

type subTypeLogic struct {
	sl subtype.Repository
}

func New(pr subtype.Repository) subtype.UseCase {
	return &subTypeLogic{
		sl: sr,
	}
}

func (stl *subTypeLogic) AddSubTypeLogic(newType subtype.Core) error {
	var queryDatas subtype.RepoData

	if newType.SubmissionTypeName == "" {
		log.Error("submission type name are empty")
		return errors.New("submission type name cannot be empty")
	}

	if len(newType.PositionTag) == 0 || len(newType.SubmissionValues) == 0 {
		log.Error("submission owners and values are empty")
		return errors.New("submission owners values cannot be empty")
	}

	

}
