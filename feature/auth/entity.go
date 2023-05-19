package auth

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
	LoginHandler() echo.HandlerFunc
	SignValidationLogic() echo.HandlerFunc
}

type UseCase interface {
	LogInLogic(id string, password string) (Core, error)
	SignVallidationLogic(signID string) (SignCore, error)
}

type Repository interface {
	Login(id string, password string) (Core, error)
	SignVaidation(signID string) (SignCore, error)
}

type SignCore struct {
	Title            string
	Officialname     string
	Officialposition string
	Date             string
}
