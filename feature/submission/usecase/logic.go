package usecase

import (
	"fmt"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
)

type submissionLogic struct {
	sl submission.Repository
}

func New(sr submission.Repository) submission.UseCase {
	return &submissionLogic{
		sl: sr,
	}
}

func (sr *submissionLogic) FindRequirementLogic(userID string, typeName string, value int) (submission.Core, error) {
	result, err := sr.sl.FindRequirement(userID, typeName, value)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return submission.Core{}, fmt.Errorf("data not found %w", err)
		} else if strings.Contains(err.Error(), "syntax") {
			return submission.Core{}, fmt.Errorf("internal server error %w", err)
		} else {
			return submission.Core{}, fmt.Errorf("unexpected error %w", err)
		}
	}

	return result, nil
}
