package user

import (
	office "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office"
	position "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"

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
	Position    position.Core
	Office      office.Core
}

type Handler interface {
	RegisterHandler() echo.HandlerFunc
	GetAllUserHandler() echo.HandlerFunc
	GetUserByIdHandler() echo.HandlerFunc
	UpdateUserHandler() echo.HandlerFunc
	DeleteUserHandler() echo.HandlerFunc
}

type UseCase interface {
	RegisterUser(newUser Core) error
	GetAllUser(limit, offset int, name string) ([]Core, error)
	GetUserById(id string) (Core, error)
	UpdateUser(id string, updateUser Core) error
	DeleteUser(id string) error
}

type Repository interface {
	InsertUser(newUser Core) error
	SelectAllUser(limit, offset int, name string) ([]Core, error)
	GetUserById(id string) (Core, error)
	UpdateUser(id string, input Core) error
	DeleteUser(id string) error
	GetPositionTagByID(positionID int) (string, error)
	GenerateIDFromPositionTag(positionTag string) (string, error)
	CheckUserIDExists(id string) (bool, error)
}
