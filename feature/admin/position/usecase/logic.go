package usecase

import (
	"errors"
	"fmt"
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

func (pl *positionLogic) AddPositionLogic(newPosition position.Core) error {
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

func (pl *positionLogic) GetPositionsLogic(limit int, offset int, search string) ([]position.Core, int64, error) {

	positions, count, err := pl.pl.GetPositions(limit, offset, search)
	if err != nil {
		log.Error("error on getpositions query")
		return nil, 0, err
	}

	return positions, count, nil
}

func (pl *positionLogic) DeletePositionLogic(position int) error {
	if err := pl.pl.DeletePosition(position); err != nil {
		if strings.Contains(err.Error(), "count query error") {
			log.Error("errors occurs when countin the datas for delete")
			return fmt.Errorf("count position query error %w", err)
		}

		if strings.Contains(err.Error(), "data found") {
			log.Error("no position data found for deletion")
			return fmt.Errorf("no data found for deletion %w", err)
		}

		log.Error("data found, but delete query error")
		return err
	}

	return nil
}
