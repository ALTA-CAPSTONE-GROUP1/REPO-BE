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

// DeleteOfficeLogic implements office.UseCase
func (ol *officeLogic) DeleteOfficeLogic(id uint) error {
	err := ol.ol.DeleteOffice(id)
	if err != nil {
		log.Error("failed on calling deleteuser query")
		if strings.Contains(err.Error(), "finding office") {
			log.Error("error on finding office (not found)")
			return errors.New("bad request, office not found")
		} else if strings.Contains(err.Error(), "cannot delete") {
			log.Error("error on delete office")
			return errors.New("internal server error, cannot delete office")
		}
		log.Error("error in delete office (else)")
		return err
	}
	return nil
}

// GetAllOfficeLogic implements office.UseCase
func (ol *officeLogic) GetAllOfficeLogic(limit int, offset int, search string) ([]office.Core, int, error) {
	result, totaldata, err := ol.ol.GetAllOffice(limit, offset, search)
	if err != nil {
		log.Error("failed to find all office", err.Error())
		return []office.Core{}, totaldata, errors.New("internal server error")
	}

	return result, totaldata, nil
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
