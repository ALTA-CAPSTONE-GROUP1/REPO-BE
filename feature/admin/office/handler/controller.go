package handler

import (
	"net/http"
	"strconv"
	"strings"

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

// DeleteOfficeHandler implements office.Handler
func (oc *officeController) DeleteOfficeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces delete office")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		officeIDStr := c.QueryParam("id")
		if officeIDStr == "" {
			c.Logger().Error("office ID is missing")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "office ID is missing", nil))
		}

		officeID, err := strconv.Atoi(officeIDStr)
		if err != nil {
			c.Logger().Error("invalid office ID format")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid office ID format", nil))
		}

		if err = oc.service.DeleteOfficeLogic(uint(officeID)); err != nil {
			c.Logger().Error("error in calling DeleteOfficeLogic")
			if strings.Contains(err.Error(), "office not found") {
				c.Logger().Error("error in calling DeleteOfficeLogic, office not found")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "office not found", nil))
			} else if strings.Contains(err.Error(), "cannot delete") {
				c.Logger().Error("error in calling DeleteOfficeLogic, cannot delete")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete office", nil))
			}

			c.Logger().Error("error in calling DeleteOfficeLogic, cannot delete")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete office", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "successfully deleted office", nil))
	}
}

// GetAllOfficeHandler implements office.Handler
func (oc *officeController) GetAllOfficeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user is not an admin trying to access get all office")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not an admin", nil))
		}

		search := c.QueryParam("search")

		limitStr := c.QueryParam("limit")
		offsetStr := c.QueryParam("offset")

		limit := 5
		if limitStr != "" {
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				c.Logger().Errorf("limit is not a number: %s", limitStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid limit value", nil))
			}
			limit = limitInt
		}

		offset := 0
		if offsetStr != "" {
			offsetInt, err := strconv.Atoi(offsetStr)
			if err != nil {
				c.Logger().Errorf("offset is not a number: %s", offsetStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid offset value", nil))
			}
			offset = offsetInt
		}

		office, totaldata, err := oc.service.GetAllOfficeLogic(limit, offset, search)
		if err != nil {
			c.Logger().Error("error occurs when calling GetAllOfficeLogic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error", nil))
		}

		pagination := helper.Pagination(limit, offset, totaldata)

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK, "Successfully retrieved office data", office, pagination))
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
			Name: req.Name,
		}

		if err := oc.service.AddOfficeLogic(newOffice); err != nil {
			c.Logger().Errorf("error occurs on calling Office Logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create office", nil))
	}
}
