package handler

import (
	"net/http"
	"strconv"

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
