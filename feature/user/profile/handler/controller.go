package handler

import (
	"net/http"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/jinzhu/copier"
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

// UpdateUserHandler implements profile.Handler
func (uc *profileController) UpdateUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		updateInput := InputUpdate{}
		if err := c.Bind(&updateInput); err != nil {
			c.Logger().Error("terjadi kesalahan bind", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		updateUser := profile.Core{}
		copier.Copy(&updateUser, &updateInput)
		if err := uc.service.UpdateUser(userID, updateUser); err != nil {
			c.Logger().Error("failed on calling updateprofile log")
			if strings.Contains(err.Error(), "hashing password") {
				c.Logger().Error("hashing password error")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "new password are invalid", nil))
			} else if strings.Contains(err.Error(), "affected") {
				c.Logger().Error("no rows affected on update user")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "data is up to date", nil))
			}

			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to update user data", nil))
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
