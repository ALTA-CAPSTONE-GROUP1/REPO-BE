package handler

import (
	"net/http"
	"strconv"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type submissionController struct {
	sc submission.UseCase
}

func New(sl submission.UseCase) submission.Handler {
	return &submissionController{
		sc: sl,
	}
}

func (sc *submissionController) FindRequirementHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var response RequirementResponseBody

		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces delete office")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		typeName := c.QueryParam("submission_type")
		value := c.QueryParam("submission_value")
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			c.Logger().Error("value cannot convert to int")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "value are cannot processed now", nil))
		}
		result, err := sc.sc.FindRequirementLogic(userID, typeName, valueInt)
		if err != nil{
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError,"server errror", nil))
		}

		
		response.To = make([]ToApprover, len(result.To))
		response.CC = make([]CcApprover, len(result.CC))

		for i, to := range result.To {
			response.To[i] = ToApprover{
				ApproverPosition: to.ApproverPosition,
				ApproverId:       to.ApproverId,
				ApproverName:     to.ApproverName,
			}
		}

		for i, cc := range result.CC {
			response.CC[i] = CcApprover{
				CcPosition: cc.CcPosition,
				CcName:     cc.CcName,
				CcId:       cc.CcId,
			}
		}

		response.Requirement = result.Requirement

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to get requirement data", response))
	}
}
