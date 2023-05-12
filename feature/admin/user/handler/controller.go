package handler

import (
	"net/http"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type userController struct {
	service user.UseCase
}

func New(u user.UseCase) user.Handler {
	return &userController{
		service: u,
	}
}

// RegisterHandler implements user.Handler
func (uc *userController) RegisterHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := RegisterInput{}
		if err := c.Bind(&input); err != nil {
			c.Logger().Error("error on bind register input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		err := uc.service.RegisterUser(user.Core{
			OfficeID:    input.OfficeID,
			PositionID:  input.PositionID,
			Name:        input.Name,
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
			Password:    input.Password,
		})
		if err != nil {
			c.Logger().Error("error on calling userLogic", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, err.Error(), nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create user", nil))
	}
}
