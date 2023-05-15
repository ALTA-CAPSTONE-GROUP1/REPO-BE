package handler

import (
	"net/http"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type profileController struct {
	service profile.UseCase
}

func New(u profile.UseCase) profile.Handler {
	return &profileController{
		service: u,
	}
}

// ProfileHandler implements profile.Handler
func (uc *profileController) ProfileHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		data, err := uc.service.ProfileLogic(userID)
		if err != nil {
			c.Logger().Error("error on calling user profile logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error", nil))
		}
		dataResponse := CoreToUserResponse(data)
		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to get user profile", dataResponse))
	}
}
