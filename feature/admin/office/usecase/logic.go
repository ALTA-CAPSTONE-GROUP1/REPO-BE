package usecase

import (
	"errors"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office"
	"github.com/labstack/gommon/log"
)

type officeLogic struct {
	ol office.Repository
}

func New(or office.Repository) office.UseCase {
	return &officeLogic{
		ol: or,
	}
}

// GetAllOfficeLogic implements office.UseCase
func (ol *officeLogic) GetAllOfficeLogic(limit int, offset int, search string) ([]office.Core, error) {
	result, err := ol.ol.GetAllOffice(limit, offset, search)
	if err != nil {
		log.Error("failed to find all office", err.Error())
		return []office.Core{}, errors.New("internal server error")
	}

	return result, nil
}

// AddOfficeLogic implements office.UseCase
func (ol *officeLogic) AddOfficeLogic(newOffice office.Core) error {
	if err := ol.ol.InsertOffice(newOffice); err != nil {
		if strings.Contains(err.Error(), "column") {
			log.Error("insert office error, COLUMN issue")
			return errors.New("server error")
		} else {
			log.Error("unexpected error when insert office")
			return err
		}
	}

	return nil
}
