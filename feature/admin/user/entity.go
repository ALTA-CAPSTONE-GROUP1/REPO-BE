package user

import (
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Core struct {
	ID          string
	OfficeID    int
	PositionID  int
	Name        string
	Email       string
	PhoneNumber string
	Password    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type Handler interface {
	RegisterHandler() echo.HandlerFunc
	GetAllUserHandler() echo.HandlerFunc
	GetUserByIdHandler() echo.HandlerFunc
	UpdateUserHandler() echo.HandlerFunc
}

type UseCase interface {
	RegisterUser(newUser Core) error
	GetAllUser(page int, name string) ([]Core, error)
	GetUserById(id string) (Core, error)
	UpdateUser(id string, updateUser Core) error
}

type Repository interface {
	InsertUser(newUser Core) error
	SelectAllUser(limit, offset int, name string) ([]Core, error)
	GetUserById(id string) (Core, error)
	UpdateUser(id string, input Core) error
}
