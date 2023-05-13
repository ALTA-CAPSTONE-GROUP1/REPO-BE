package subtype

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/labstack/echo/v4"
)

type Core struct {
	Name        string
	Requirement string
	Positions   []admin.Position
}

type Handler interface {
	AddTypeHandler() echo.HandlerFunc
}

type UseCase interface {
	AddSubTypeLogic(newType Core) error
}

type Repository interface {
	InsertSubType(newType Core) error
}
