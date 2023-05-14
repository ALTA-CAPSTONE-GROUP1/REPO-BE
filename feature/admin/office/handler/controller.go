package handler

import (
	"net/http"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type officeController struct {
	service office.UseCase
}

func New(ou office.UseCase) office.Handler {
	return &officeController{
		service: ou,
	}
}

// AddOfficeHandler implements office.Handler
func (oc *officeController) AddOfficeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(AddOfficeRequest)
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add office")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		if err := c.Bind(req); err != nil {
			c.Logger().Error("error on binding user input")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		newOffice := office.Core{
			Name: req.Office,
		}

		if err := oc.service.AddOfficeLogic(newOffice); err != nil {
			c.Logger().Errorf("error occurs on calling Office Logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create office", nil))
	}
}
