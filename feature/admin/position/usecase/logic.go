package usecase

import (
	"errors"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
	"github.com/labstack/gommon/log"
)

type positionLogic struct {
	pl position.Repository
}

func New(pr position.Repository) position.UseCase {
	return &positionLogic{
		pl: pr,
	}
}

func (pl positionLogic) AddPositionLogic(newPosition position.Core) error {
	if err := pl.pl.InsertPosition(newPosition); err != nil {
		if strings.Contains(err.Error(), "column") {
			log.Error("insert position error, COLUMN issue")
			return errors.New("server error")
		} else {
			log.Error("unexpected error when insert position")
			return err
		}
	}

	return nil
}
