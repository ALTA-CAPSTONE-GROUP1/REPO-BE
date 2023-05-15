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

		officePath, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		if err = oc.service.DeleteOfficeLogic(uint(officePath)); err != nil {
			c.Logger().Error("error in calling DeletOfficeLogic")
			if strings.Contains(err.Error(), "office not found") {
				c.Logger().Error("error in calling DeletOfficeLogic, office not found")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "office not found", nil))

			} else if strings.Contains(err.Error(), "cannot delete") {
				c.Logger().Error("error in calling DeletOfficeLogic, cannot delete")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete office", nil))
			}

			c.Logger().Error("error in calling DeletOfficeLogic, cannot delete")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete office", nil))

		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to delete office", nil))
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

		limitStr := c.QueryParam("limit")
		offsetStr := c.QueryParam("offset")
		search := c.QueryParam("search")

		limit := -1
		if limitStr != "" {
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				c.Logger().Errorf("limit is not a number: %s", limitStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid limit value", nil))
			}
			limit = limitInt
		}

		offset := -1
		if offsetStr != "" {
			offsetInt, err := strconv.Atoi(offsetStr)
			if err != nil {
				c.Logger().Errorf("offset is not a number: %s", offsetStr)
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid offset value", nil))
			}
			offset = offsetInt
		}

		if limit < 0 || offset < 0 {
			c.Logger().Error("error occurs because limit or offset is negative")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Limit and offset cannot be negative", nil))
		}

		office, err := oc.service.GetAllOfficeLogic(limit, offset, search)
		if err != nil {
			c.Logger().Error("error occurs when calling GetAllOfficeLogic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "Successfully retrieved office data", office))
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
			Name:     req.Name,
			Level:    req.Level,
			ParentID: req.ParentID,
		}

		if err := oc.service.AddOfficeLogic(newOffice); err != nil {
			c.Logger().Errorf("error occurs on calling Office Logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create office", nil))
	}
}
