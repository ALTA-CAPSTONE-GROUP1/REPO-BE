package position

import (
	"github.com/labstack/echo/v4"
)

type Core struct {
	ID   int
	Name string
	Tag  string
}

type Handler interface {
	AddPositionHandler() echo.HandlerFunc
	GetAllPositionHandler() echo.HandlerFunc
	DeletePositionHandler() echo.HandlerFunc
}

type UseCase interface {
	AddPositionLogic(newPosition Core) error
	GetPositionsLogic(limit int, offset int, search string) ([]Core, int64, error)
	DeletePositionLogic(position int) error
}

type Repository interface {
	InsertPosition(newPosition Core) error
	GetPositions(limit int, offset int, search string) ([]Core, int64, error)
	DeletePosition(position int) error
}
