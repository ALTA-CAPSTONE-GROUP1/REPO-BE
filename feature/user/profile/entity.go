package profile

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/labstack/echo/v4"
)

type Core struct {
	ID          string
	OfficeID    int
	PositionID  int
	Name        string
	Email       string
	PhoneNumber string
	Password    string
	Position    admin.Position
	Office      admin.Office
}
type Handler interface {
	ProfileHandler() echo.HandlerFunc
	UpdateUserHandler() echo.HandlerFunc
}

type UseCase interface {
	ProfileLogic(id string) (Core, error)
	UpdateUser(id string, updateUser Core) error
}

type Repository interface {
	Profile(id string) (Core, error)
	UpdateUser(id string, input Core) error
}
