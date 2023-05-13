package handler

import (
	"net/http"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type subTypeController struct {
	service subtype.UseCase
}

func New(su subtype.UseCase) subtype.Handler {
	return &subTypeController{
		service: su,
	}
}

func (sc *subTypeController) AddTypeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(AddSubTypeRequest)
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		if err := c.Bind(req); err != nil {
			c.Logger().Error("error on binding user input")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		if err := c.Validate(req); err != nil {
			c.Logger().Error("errror in validate input" + err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "bad request, invalid input", nil))
		}

		newSubType := subtype.Core{}
		newSubType.SubmissionTypeName = req.SubmissionTypeName
		newSubType.PositionTag = req.Position

		for _, subVal := range req.SubmissionValues {
			newSubVal := subtype.ValueDetails{
				Value:         subVal.Value,
				TagPositionTo: subVal.PositionTo,
				TagPositionCC: subVal.PositionCC,
			}
			newSubType.SubmissionValues = append(newSubType.SubmissionValues, newSubVal)
		}

		newSubType.Requirement = req.Requirement

		if err := sc.service.AddSubTypeLogic(newSubType); err != nil {
			if strings.Contains(err.Error(), "failed to insert submission type data") {
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to insert submission type data", nil))
			} else if strings.Contains(err.Error(), "failed to add user as authorized to make this submission type") {
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Failed to add user as authorized to make this submission type", nil))
			} else if strings.Contains(err.Error(), "cannot find authorized officials approver by tag") {
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Cannot find authorized officials approver by tag", nil))
			} else if strings.Contains(err.Error(), "cannot find authorized officials ccs by tag") {
				return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Cannot find authorized officials ccs by tag", nil))
			} else if strings.Contains(err.Error(), "failed to add roles to data type") {
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to add roles to data type", nil))
			} else if strings.Contains(err.Error(), "failed to save all data to database") {
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to save all data to database", nil))
			} else {
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Unexpected error", nil))
			}
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to create submission type", nil))
	}
}
