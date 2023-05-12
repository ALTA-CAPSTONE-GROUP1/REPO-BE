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
	GetAllPositionHandler() echo.HandlerFunc
}

type UseCase interface {
	AddPositionLogic(newPosition Core) error
	GetPositionsLogic(limit int, offset int, search string) ([]Core, error)
}

type Repository interface {
	InsertPosition(newPosition Core) error
	GetPositions(limit int, offset int, search string) ([]Core, error)
}
