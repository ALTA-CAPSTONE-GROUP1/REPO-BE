package position

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Core struct {
	Name      string
	Tag       string
	Types     []admin.Type
	DeletedAt gorm.DeletedAt
}

type Handler interface {
	AddPositionHandler() echo.HandlerFunc
}

type UseCase interface {
	AddPositionLogic(newPosition Core) error
}

type Repository interface {
	InsertPositionHandler(newPosition Core) error
}
