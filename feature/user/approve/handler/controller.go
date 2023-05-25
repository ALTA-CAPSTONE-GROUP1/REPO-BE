package handler

import (
	"net/http"
	"strconv"
	"strings"

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

// UpdateSubmissionApproveHandler implements approve.Handler
func (ac *approveController) UpdateSubmissionApproveHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		submissionID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		updateInput := InputUpdate{}
		if err := c.Bind(&updateInput); err != nil {
			c.Logger().Error("terjadi kesalahan bind", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		input := approve.Core{
			Status:  updateInput.Action,
			Message: updateInput.Message,
		}

		if err := ac.service.UpdateApprove(userID, submissionID, input); err != nil {
			c.Logger().Error("failed on calling updateprofile log")
			if strings.Contains(err.Error(), "affected") {
				c.Logger().Error("no rows affected on update submission")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "data is up to date", nil))
			}

			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to give action", nil))
	}
}

// GetSubmissionByIdHandler implements approve.Handler
func (ac *approveController) GetSubmissionByIdHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == "" {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		id, errid := strconv.Atoi(c.Param("id"))
		if errid != nil {
			c.Logger().Error("cannot use path param", errid.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		data, err := ac.service.GetSubmissionById(userID, id)
		if err != nil {
			c.Logger().Error("error on calling get all user logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to read data", nil))
		}

		dataResponse := CoreToApproveByIdResponse(data)
		return c.JSON(helper.ResponseFormat(http.StatusOK, "Successfully retrieved all users", dataResponse))
	}
}

// GetSubmissionAprroveHandler implements approve.Handler
func (ac *approveController) GetSubmissionAprroveHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		limitStr := c.QueryParam("limit")
		offsetStr := c.QueryParam("offset")
		title := c.QueryParam("title")
		fromParam := c.QueryParam("from")
		types := c.QueryParam("type")

		limit := 10
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

		search := approve.GetAllQueryParams{
			Title:  title,
			FromTo: fromParam,
			Type:   types,
			Limit:  limit,
			Offset: offset,
		}

		data, totalData, err := ac.service.GetSubmissionAprrove(userID, search)
		if err != nil {
			c.Logger().Error("error on calling get all user logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to read data", nil))
		}

		pagination := helper.Pagination(limit, offset, int(totalData))
		dataResponse := CoreToGetAllApproveResponse(data, userID)
		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK, "Successfully retrieved all submissions", dataResponse, pagination))
	}
}
