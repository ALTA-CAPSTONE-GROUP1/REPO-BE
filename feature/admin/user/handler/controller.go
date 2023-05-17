package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/jinzhu/copier"
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

// DeleteUserHandler implements user.Handler
func (uc *userController) DeleteUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces delete user account")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		userPath := c.Param("id")

		err := uc.service.DeleteUser(userPath)
		if err != nil {
			c.Logger().Error("error in calling DeletUserLogic")
			if strings.Contains(err.Error(), "user not found") {
				c.Logger().Error("error in calling DeletUserLogic, user not found")
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "user not found", nil))

			} else if strings.Contains(err.Error(), "cannot delete") {
				c.Logger().Error("error in calling DeletUserLogic, cannot delete")
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete user", nil))
			}

			c.Logger().Error("error in calling DeletUserLogic, cannot delete")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error in delete user", nil))

		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to delete user", nil))
	}
}

// UpdateUserHandler implements user.Handler
func (uc *userController) UpdateUserHandler() echo.HandlerFunc {
	return func(c echo.Context) error {

		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}
		userPath := c.Param("id")

		updateInput := InputUpdate{}
		if err := c.Bind(&updateInput); err != nil {
			c.Logger().Error("terjadi kesalahan bind", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		updateUser := user.Core{}
		copier.Copy(&updateUser, &updateInput)
		if err := uc.service.UpdateUser(userPath, updateUser); err != nil {
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
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		userPath := c.Param("id")
		// if err != nil {
		// 	c.Logger().Error("cannot use path param", err.Error())
		// 	return c.JSON(helper.ResponseFormat(http.StatusNotFound, "path invalid", nil))
		// }

		data, err := uc.service.GetUserById(userPath)
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
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user is not an admin trying to access get all office")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not an admin", nil))
		}

		limitStr := c.QueryParam("limit")
		offsetStr := c.QueryParam("offset")
		name := c.QueryParam("name")

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

		data, err := uc.service.GetAllUser(limit, offset, name)
		if err != nil {
			c.Logger().Error("error on calling get all user logic")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to read data", nil))
		}

		dataResponse := CoreToGetAllUserResponse(data)

		pagination := helper.Pagination(limit, offset, len(data))

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK, "Successfully retrieved all users", dataResponse, pagination))
	}
}

// RegisterHandler implements user.Handler
func (uc *userController) RegisterHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user is not admin trying to access add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not an admin", nil))
		}

		input := RegisterInput{}
		if err := c.Bind(&input); err != nil {
			c.Logger().Error("error on binding register input", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		err := uc.service.RegisterUser(user.Core{
			OfficeID:    input.OfficeID,
			PositionID:  input.PositionID,
			Name:        input.Name,
			Email:       input.Email,
			PhoneNumber: input.PhoneNumber,
		})
		if err != nil {
			c.Logger().Error("error on calling RegisterUser", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, err.Error(), nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "successfully created user", nil))
	}
}
