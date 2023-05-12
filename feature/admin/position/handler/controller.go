package handler

import (
	"net/http"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
			log.Error("user are not admin try to acces add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}
		
		if err := c.Bind(req); err != nil {
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		if err := c.Validate(req); err != nil {
			log.Error("errror in validate input" + err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "bad request, invalid input", nil))
		}
		newPosition := position.Core{
			Name: req.Position,
			Tag:  req.Tag,
		}

		if err := pc.service.AddPositionLogic(newPosition); err != nil {
			log.Errorf("error occurs on calling Position logic with data %v, %v", newPosition.Name, newPosition.Tag)
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create position", nil))
	}
}
