package handler

import (
	"math"
	"net/http"
	"strconv"
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

		if err := c.Bind(&req); err != nil {
			c.Logger().Error("error on binding user input")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "invalid input", nil))
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

		return c.JSON(helper.ResponseFormat(http.StatusCreated, "succes to create submission type", nil))
	}
}

func (sc *subTypeController) GetTypesHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		res := new(GetSubmissionTypeResponse)
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		limit := c.QueryParam("limit")
		offset := c.QueryParam("offset")
		search := c.QueryParam("search")

		if limit == "" {
			limit = "10"
		}

		if offset == "" {
			offset = "0"
		}

		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			c.Logger().Errorf("limit are not a number %v", limit)
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "Server Error, limit are NaN", nil))
		}
		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			c.Logger().Errorf("offset are not a number %v", offset)
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Server Error, offset are NaN", nil))
		}
		if limitInt < 0 || offsetInt < 0 {
			c.Logger().Error("error occurs because limit/offset are negatif")
			return c.JSON(helper.ResponseFormat(http.StatusBadRequest, "limit and offset cannot negative", nil))
		}

		subTypesData, positionsData, err := sc.service.GetSubTypesLogic(limitInt, offsetInt, search)
		if err != nil {
			if strings.Contains(err.Error(), "retrieve positions") {
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to retrieve positions", nil))
			}
			if strings.Contains(err.Error(), "retrieve submission types") {
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to retrieve submission types", nil))
			}
			if strings.Contains(err.Error(), "retrieve position_has_types") {
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to retrieve position_has_types", nil))
			}
			c.Logger().Error("unexpected error on calling all submisioon type data")
			return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "Failed to get submission types with unexpected error", nil))
		}

		for _, p := range positionsData {
			position := Position{
				PositionName: p.PositionName,
				PositionTag:  p.PositionTag,
			}
			res.Position = append(res.Position, position)
		}

		filteredData := []subtype.GetSubmissionTypeCore{}
		if search != "" {
			for _, data := range subTypesData {
				if strings.Contains(strings.ToLower(data.SubmissionTypeName), strings.ToLower(search)) ||
					strings.Contains(strings.ToLower(data.Requirement), strings.ToLower(search)) {
					filteredData = append(filteredData, data)
				}
			}
		} else {
			filteredData = subTypesData
		}

		for _, std := range filteredData {
			found := false
			for i, r := range res.SubmissionType {
				if r.SubmissionTypeName == std.SubmissionTypeName {
					tmpsd := SubmissionDetail{
						SubmissionValue:       std.Value,
						SubmissionRequirement: std.Requirement,
					}
					tmpSubDetail := []SubmissionDetail{}
					for _, sd := range r.SubmissionDetail {
						if sd.SubmissionValue == tmpsd.SubmissionValue && sd.SubmissionRequirement == tmpsd.SubmissionRequirement {
							found = true
							break
						}
						tmpSubDetail = append(tmpSubDetail, sd)
					}
					tmpSubDetail = append(tmpSubDetail, tmpsd)
					res.SubmissionType[i].SubmissionDetail = tmpSubDetail
					found = true
					break
				}
			}

			if !found {
				newSubmissionType := SubmissionType{
					SubmissionTypeName: std.SubmissionTypeName,
					SubmissionDetail: []SubmissionDetail{
						{
							SubmissionValue:       std.Value,
							SubmissionRequirement: std.Requirement,
						},
					},
				}
				res.SubmissionType = append(res.SubmissionType, newSubmissionType)
			}
		}
		if offsetInt+limitInt > len(filteredData) {
			limitInt = len(filteredData) - offsetInt
		}
		if offsetInt+limitInt> len(res.SubmissionType){
			limitInt = len(res.SubmissionType) - offsetInt
		}
		res.SubmissionType = res.SubmissionType[offsetInt : offsetInt+limitInt]

		totalData := len(filteredData)
		totalPage := int(math.Ceil(float64(totalData) / float64(limitInt)))
		currentPage := int(math.Ceil(float64(offsetInt+1) / float64(limitInt)))

		meta := Meta{
			CurrentLimit:  limitInt,
			CurrentOffset: offsetInt,
			CurrentPage:   currentPage,
			TotalData:     totalData,
			TotalPage:     totalPage,
		}

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK, "succes to get submission types data", res, meta))
	}
}

func (sc *subTypeController) DeleteTypeHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID != "admin" {
			c.Logger().Error("user are not admin try to acces add position")
			return c.JSON(helper.ResponseFormat(http.StatusUnauthorized, "you are not admin", nil))
		}

		subTypeName := c.QueryParam("submission_type")

		if err := sc.service.DeleteSubTypeLogic(subTypeName); err != nil {
			c.Logger().Errorf("error on delete subtype", err)
			if strings.Contains(err.Error(), "empty set") {
				return c.JSON(helper.ResponseFormat(http.StatusNotFound, "subtype not found please refresh and try again", nil))
			} else {
				return c.JSON(helper.ResponseFormat(http.StatusInternalServerError, "server error", nil))
			}
		}

		return c.JSON(helper.ResponseFormat(http.StatusOK, "succes to delete submission type data", nil))
	}

}
