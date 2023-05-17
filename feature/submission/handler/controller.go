package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
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
		if userID == "" {
			c.Logger().Error("")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "invalid or expired JWT", nil))
		}

		typeName := c.QueryParam("submission_type")
		value := c.QueryParam("submission_value")
		valueInt, err := strconv.Atoi(value)
		if err != nil {
			c.Logger().Error("value cannot convert to int")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "value are cannot processed now", nil))
		}
		result, err := sc.sc.FindRequirementLogic(userID, typeName, valueInt)
		if err != nil {
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server errror", nil))
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

func (sc *submissionController) AddSubmissionHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		var newSub submission.AddSubmissionCore
		userID := helper.DecodeToken(c)
		if userID == "" {
			c.Logger().Error("invalid or expired jwt")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "invalid or expired JWT", nil))
		}

		req := new(AddAddSubReq)
		if err := c.Bind(req); err != nil {
			log.Errorf("error on finding binding submission", err)
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest,
				"bad request",
				nil))
		}

		attachmentHeader, err := c.FormFile("attachment")
		if err != nil {
			log.Error("error occurs on read attachment")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest,
				"bad request",
				nil,
			))
		}

		if err != nil {
			log.Error("error occurs on open attachment")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "internal server error", nil))
		}

		newSub.Message = req.Message
		newSub.SubmissionType = req.SubmissionType
		newSub.SubmissionValue = req.SubmissionValue
		newSub.OwnerID = userID
		newSub.Title = req.Title

		for _, v := range req.CC {
			tmp := submission.CcApprover{
				CcId: v,
			}
			newSub.CC = append(newSub.CC, tmp)
		}

		for _, v := range req.To {
			tmp := submission.ToApprover{
				ApproverId: v,
			}
			newSub.ToApprover = append(newSub.ToApprover, tmp)
		}

		if err := sc.sc.AddSubmissionLogic(newSub, attachmentHeader); err != nil {
			log.Error("error on calling addsubmissionlogic")
			if strings.Contains(err.Error(), "record not found ") {
				return c.JSON(helper.ResponseFormat(
					http.StatusNotFound,
					"data not found",
					nil,
				))
			}
			if strings.Contains(err.Error(), "syntax") {
				return c.JSON(
					helper.ResponseFormat(http.StatusInternalServerError,
						"internal server error",
						nil))
			}
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error", nil))
		}

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create submission", nil))
	}
}

func (sc *submissionController) GetAllSubmissionHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == "" {
			c.Logger().Error("invalid or expired jwt")
			return c.JSON(helper.ReponseFormatWithMeta(http.StatusUnauthorized, "invalid or expired JWT", nil, nil))
		}

		var params submission.GetAllQueryParams
		limit := c.QueryParam("limit")
		offset := c.QueryParam("offset")
		if limit == "" {
			limit = "10"
		}

		if offset == "" {
			offset = "0"
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			c.Logger().Error("cannot convert limit to int")
			return c.JSON(helper.ReponseFormatWithMeta(http.StatusBadRequest,
				"limit must be string",
				nil, nil))
		}
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			c.Logger().Error("cannot convert offset to int")
			return c.JSON(helper.ReponseFormatWithMeta(http.StatusBadRequest,
				"offset must be string",
				nil, nil))
		}
		searchInTitle := c.QueryParam("title")
		searchInTo := c.QueryParam("to")

		submissionDatas, subTypeDatas, err := sc.sc.GetAllSubmissionLogic(userID, params)
		if err != nil {
			if strings.Contains(err.Error(), "record") {
				return c.JSON(helper.ReponseFormatWithMeta(http.StatusNotFound, "record not found", nil, nil))
			}
			return c.JSON(helper.ReponseFormatWithMeta(http.StatusInternalServerError, "Server error", nil, nil))
		}

		var submissions []Submission
		filteredData := []submission.AllSubmiisionCore{}

		if searchInTitle != "" {
			for _, data := range submissionDatas {
				if strings.Contains(strings.ToLower(data.Title), strings.ToLower(searchInTitle)) {
					filteredData = append(filteredData, data)
				}
			}
		} else {
			filteredData = submissionDatas
		}

		if searchInTo != "" {
			for _, data := range filteredData {
				for _, to := range data.Tos {
					if strings.Contains(strings.ToLower(to.ApproverName), strings.ToLower(searchInTo)) ||
						strings.Contains(strings.ToLower(to.ApproverId), strings.ToLower(searchInTo)) ||
						strings.Contains(strings.ToLower(to.ApproverPosition), strings.ToLower(searchInTo)) {
						filteredData = append(filteredData, data)
					}
				}
			}
		}

		for _, submissionData := range filteredData {
			var toApprovers []Approver
			for _, to := range submissionData.Tos {
				toApprovers = append(toApprovers, Approver{
					ApproverPosition: to.ApproverPosition,
					ApproverName:     to.ApproverName,
				})
			}
			var ccApprovers []CC
			for _, cc := range submissionData.CCs {
				ccApprovers = append(ccApprovers, CC{
					CCPosition: cc.CcPosition,
					CCName:     cc.CcName,
				})
			}

			submissions = append(submissions, Submission{
				ID:             submissionData.ID,
				To:             toApprovers,
				CC:             ccApprovers,
				Title:          submissionData.Title,
				Status:         submissionData.Status,
				Attachment:     submissionData.Attachment,
				ReceiveDate:    submissionData.ReceiveDate,
				Opened:         submissionData.Opened,
				SubmissionType: submissionData.SubmissionType,
			})
		}

		var submissionTypeChoices []SubmissionTypeChoice
		for _, subTypeData := range subTypeDatas {
			var values []int
			for _, pos := range subTypeData.Positions {
				values = append(values, pos.ID)
			}

			submissionTypeChoices = append(submissionTypeChoices, SubmissionTypeChoice{
				Name:   subTypeData.Name,
				Values: values,
			})
		}

		if offsetInt+limitInt > len(filteredData) {
			limitInt = len(filteredData) - offsetInt
		}
		submissions = submissions[offsetInt : offsetInt+limitInt]
		var totalData int
		var totalPage int
		if len(filteredData) > 0 {
			totalData = len(filteredData)
			totalPage = int(math.Ceil(float64(totalData) / float64(limitInt)))
		}
		currentPage := int(math.Ceil(float64(offsetInt+1) / float64(limitInt)))

		meta := Meta{
			CurrentLimit:  limitInt,
			CurrentOffset: offsetInt,
			CurrentPage:   currentPage,
			TotalData:     totalData,
			TotalPage:     totalPage,
		}

		response := SubmissionResponse{
			Submissions:           submissions,
			SubmissionTypeChoices: submissionTypeChoices,
		}

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK, "succes to get submissions data", response, meta))
	}
}

func (sc *submissionController) GetSubmissionByIdHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == "" {
			c.Logger().Error("invalid or expired jwt")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "invalid or expired JWT", nil))
		}
		IDParam := c.Param("submisision_id")
		subID, err := strconv.Atoi(IDParam)
		if err != nil {
			log.Errorf("error on convert submissionID to int", err.Error())
			return c.JSON(helper.ResponseFormat(
				http.StatusBadRequest,
				"Bad Request, subID must be a number",
				nil,
			))
		}

		result, err := sc.sc.GetSubmissionByIDLogic(subID)
		if err != nil {
			log.Errorf("error in calling submissionID Logic", err)
			if strings.Contains(err.Error(), "syntax") {
				return c.JSON(helper.ResponseFormat(
					http.StatusInternalServerError,
					"server error",
					nil,
				))
			}
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(helper.ResponseFormat(
					http.StatusNotFound,
					"submission not found",
					nil,
				))
			}

			return c.JSON(helper.ResponseFormat(
				http.StatusInternalServerError,
				"server error",
				nil,
			))
		}

		var response ResponseByID
		for _, to := range result.ApproverActions {
			tmp := ApproverRecipient{
				ApproverPosition: to.ApproverPosition,
				ApproverName:     to.ApproverName,
			}
			tmpAction := ApproverAction{
				ApproverName:     to.ApproverName,
				ApproverPosition: to.ApproverPosition,
				Action:           to.Action,
			}
			response.To = append(response.To, tmp)
			response.ApproverAction = append(response.ApproverAction, tmpAction)
		}

		for _, cc := range result.CC {
			tmp := CCRecipient{
				CCPosition: cc.CcPosition,
				CCName:     cc.CcName,
			}
			response.CC = append(response.CC, tmp)
		}

		response.Attachment = result.Attachment
		response.Title = result.Title
		response.ActionMessage = result.ActionMessage
		response.Message = result.Message
		response.SubmissionType = result.SubmissionType

		return c.JSON(helper.ResponseFormat(
			http.StatusOK,
			"succes to get submission by id",
			response,
		))
	}
}
