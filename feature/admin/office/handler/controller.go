package handler

import (
	"net/http"
	"strconv"

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

// GetAllOfficeHandler implements office.Handler
func (oc *officeController) GetAllOfficeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces get all office")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		limit := c.QueryParam("limit")
		offset := c.QueryParam("offset")
		search := c.QueryParam("search")

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			c.Logger().Errorf("limit are not a number %v", limit)
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Server Error, limit are NaN", nil))
		}
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			c.Logger().Errorf("offset are not a number %v", offset)
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error, offset are NaN", nil))
		}
		if limitInt < 0 || offsetInt < 0 {
			c.Logger().Error("error occurs because limit/offset are negatif")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "limit and offset cannot negative", nil))
		}

		office, err := oc.service.GetAllOfficeLogic(limitInt, offsetInt, search)
		if err != nil {
			c.Logger().Error("error occurs when calling GetOfficeLogic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to get office data", office))
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
