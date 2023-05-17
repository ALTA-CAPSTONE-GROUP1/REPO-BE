package handler

import (
	"net/http"
	"strconv"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type approveController struct {
	service approve.UseCase
}

func New(a approve.UseCase) approve.Handler {
	return &approveController{
		service: a,
	}
}

// GetSubmissionAprroveHandler implements approve.Handler
func (ac *approveController) GetSubmissionAprroveHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
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

		data, err := ac.service.GetSubmissionAprrove(userID, limit, offset, search)
		if err != nil {
			c.Logger().Error("error on calling get all user logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to read data", nil))
		}

		dataResponse := CoreToGetAllApproveResponse(data)
		return c.JSON(helper.ResponseFormat(http.StatusOK, "Successfully retrieved all users", dataResponse))
	}
}
