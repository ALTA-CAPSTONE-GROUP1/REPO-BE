package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
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

// GetSubmissionByHyperApprovalHandler implements approve.Handler
func (ac *approveController) GetSubmissionByHyperApprovalHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user is not admin trying to access delete user account")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "You are not an admin", nil))
		}

		submissionID := c.FormValue("submission_id")
		if submissionID == "" {
			c.Logger().Error("submission ID is missing")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Submission ID is required", nil))
		}

		token := c.FormValue("token")
		if token == "" {
			c.Logger().Error("token is missing")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Token is required", nil))
		}

		submissionIDInt, err := strconv.Atoi(submissionID)
		if err != nil {
			c.Logger().Error("invalid submission ID")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Invalid submission ID", nil))
		}

		result, err := ac.service.GetSubmissionByHyperApproval(userID, submissionIDInt, token)
		if err != nil {
			c.Logger().Error("error on calling get submission by hyper approval logic:", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to read data", nil))
		}

		var response ResponseByID
		for _, to := range result.ApproverActions {
			tmpAction := ApproverAction{
				ApproverName:     to.ApproverName,
				ApproverPosition: to.ApproverPosition,
				Action:           to.Action,
			}
			response.ApproverAction = append(response.ApproverAction, tmpAction)
		}

		response.Attachment = result.Attachment
		response.Title = result.Title
		response.Message = result.Message
		response.SubmissionType = result.SubmissionType

		return c.JSON(helper.ResponseFormat(
			http.StatusOK,
			"succes to get submission by id",
			response,
		))
	}
}

// UpdateSubmissionApproveHandler implements approve.Handler
func (ac *approveController) UpdateByHyperApprovalHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces delete user account")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		updateInput := InputUpdate{}
		if err := c.Bind(&updateInput); err != nil {
			c.Logger().Error("terjadi kesalahan bind", err.Error())
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
		}

		input := approve.Core{
			ID:     updateInput.SubID,
			Status: updateInput.Action,
		}

		if err := ac.service.UpdateByHyperApproval(userID, input); err != nil {
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
