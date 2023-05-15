package office

import (
	"github.com/labstack/echo/v4"
)

type Core struct {
	ID       uint
	Name     string
	Level    string
	ParentID uint
	Parent   *Core
}

type Handler interface {
	AddOfficeHandler() echo.HandlerFunc
	GetAllOfficeHandler() echo.HandlerFunc
	DeleteOfficeHandler() echo.HandlerFunc
}

type UseCase interface {
	AddOfficeLogic(newOffice Core) error
	GetAllOfficeLogic(limit int, offset int, search string) ([]Core, error)
	DeleteOfficeLogic(id uint) error
}

type Repository interface {
	InsertOffice(newOffice Core) error
	GetAllOffice(limit int, offset int, search string) ([]Core, error)
	DeleteOffice(id uint) error
}
