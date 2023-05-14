package office

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Core struct {
	Name      string
	DeletedAt gorm.DeletedAt
}

type Handler interface {
	AddOfficeHandler() echo.HandlerFunc
	GetAllOfficeHandler() echo.HandlerFunc
}

type UseCase interface {
	AddOfficeLogic(newOffice Core) error
	GetAllOfficeLogic(limit int, offset int, search string) ([]Core, error)
}

type Repository interface {
	InsertOffice(newOffice Core) error
	GetAllOffice(limit int, offset int, search string) ([]Core, error)
}
