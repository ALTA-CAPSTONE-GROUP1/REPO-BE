package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type positionController struct {
	service position.UseCase
}

func New(pu position.UseCase) position.Handler {
	return &positionController{
		service: pu,
	}
}

func (pc *positionController) AddPositionHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(AddPositionRequest)
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		if err := c.Bind(req); err != nil {
			c.Logger().Error("error on binding user input")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}
		newPosition := position.Core{
			Name: req.Position,
			Tag:  req.Tag,
		}

		if err := pc.service.AddPositionLogic(newPosition); err != nil {
			c.Logger().Errorf("error occurs on calling Position logic with data %v, %v", newPosition.Name, newPosition.Tag)
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create position", nil))
	}
}

func (pc *positionController) GetAllPositionHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add position")
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

		positions, err := pc.service.GetPositionsLogic(limitInt, offsetInt, search)
		if err != nil {
			c.Logger().Error("error occurs when calling GetPositionsLogic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error", nil))
		}
		fmt.Println(positions)
		response := []GetAllPositionResponse{}

		for _, v := range positions {
			tmp := GetAllPositionResponse{
				PositionID: v.ID,
				Position:   v.Name,
				Tag:        v.Tag,
			}
			response = append(response, tmp)
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to get positions data", response))
	}
}

func (pc *positionController) DeletePositionHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}
		position := c.QueryParam("position_id")

		if position == "" {
			c.Logger().Error("position or tag are empty string")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "data to delete are empty", nil))
		}

		positionINT, err := strconv.Atoi(position)
		if err != nil {
			c.Logger().Error("position id is not a number")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "position_id must be a number", nil))
		}
		if err := pc.service.DeletePositionLogic(positionINT); err != nil {
			if strings.Contains(err.Error(), "count position query error") {
				c.Logger().Error("errors occurs when counting the datas for delete")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error", nil))
			}
			if strings.Contains(err.Error(), "no data found for deletion") {
				c.Logger().Error("no position data found for deletion")
				return c.JSON(helper.ResponseFormat(http.StatusNotFound, "position data not found", nil))
			}
			c.Logger().Error("unexpected error")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "unexpected server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to delete position data", nil))
	}
}
