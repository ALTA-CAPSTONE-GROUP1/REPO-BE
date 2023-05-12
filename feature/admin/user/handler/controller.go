package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type userController struct {
	service user.UseCase
}

func New(u user.UseCase) user.Handler {
	return &userController{
		service: u,
	}
}

// UpdateUserHandler implements user.Handler
func (uc *userController) UpdateUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		updateInput := InputUpdate{}
		userId := helper.DecodeToken(c)
		if userId == "" {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}
		userPath, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		if userId != string(userPath) {
			c.Logger().Error("userpath is not equal with userId")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "user are not authorized to delete other user account", nil))
		}

		updateUser := user.Core{
			ID:          updateInput.ID,
			OfficeID:    updateInput.OfficeID,
			PositionID:  updateInput.PositionID,
			Name:        updateInput.Name,
			Email:       updateInput.Email,
			PhoneNumber: updateInput.PhoneNumber,
			Password:    updateInput.Password,
		}

		if err := uc.service.UpdateUser(userId, updateUser); err != nil {
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

// GetUserByIdHandler implements user.Handler
func (uc *userController) GetUserByIdHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := helper.DecodeToken(c)
		if userId == "" {
			c.Logger().Error("decode token is blank")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "jwt invalid", nil))
		}

		userPath, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.Logger().Error("cannot use path param", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		}

		data, err := uc.service.GetUserById(string(userPath))
		if err != nil {
			c.Logger().Error("error on calling user by id logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error", nil))
		}
		dataResponse := CoreToUserResponse(data)
		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to get user by id", dataResponse))
	}
}

// GetAllUserHandler implements user.Handler
func (uc *userController) GetAllUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var pageNumber int = 1
		pageParam := c.QueryParam("page")
		if pageParam != "" {
			pageConv, errConv := strconv.Atoi(pageParam)
			if errConv != nil {
				c.Logger().Error("cannot read data")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed, page must number", nil))
			} else {
				pageNumber = pageConv
			}
		}

		nameParam := c.QueryParam("name")
		data, err := uc.service.GetAllUser(pageNumber, nameParam)
		if err != nil {
			c.Logger().Error("error on calling get all user logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed, error read data", nil))
		}
		dataResponse := CoreToGetAllUserResponse(data)
		return c.JSON(helper.ResponseFormat(http.StatusCreated, "get all user successfully", dataResponse))
	}
}

// RegisterHandler implements user.Handler
func (uc *userController) RegisterHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := RegisterInput{}
		if err := c.Bind(&input); err != nil {
			c.Logger().Error("error on bind register input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		err := uc.service.RegisterUser(user.Core{
			OfficeID:    input.OfficeID,
			PositionID:  input.PositionID,
			Name:        input.Name,
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
			Password:    input.Password,
		})
		if err != nil {
			c.Logger().Error("error on calling userLogic", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, err.Error(), nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create user", nil))
	}
}
