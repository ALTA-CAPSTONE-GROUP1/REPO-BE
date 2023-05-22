package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
)

type ccController struct {
	ch cc.UseCase
}

func New(cl cc.UseCase) cc.Handler {
	return &ccController{
		ch: cl,
	}
}

func (ch *ccController) GetAllCcHander() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID := helper.DecodeToken(c)
		if userID == "" {
			c.Logger().Error("")
			return c.JSON(helper.ReponseFormatWithMeta(http.StatusUnauthorized, "invalid or expired JWT", nil, nil))
		}

		searchInTitle := c.QueryParam("title")
		searchInTo := c.QueryParam("to")
		searchInFrom := c.QueryParam("from")
		searchInType := c.QueryParam("type")
		limit := c.QueryParam("limit")
		offset := c.QueryParam("offset")

		var response []Response

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

		ccDatas, err := ch.ch.GetAllCcLogic(userID)
		if err != nil {
			if strings.Contains(err.Error(), "record") {
				return c.JSON(helper.ReponseFormatWithMeta(http.StatusInternalServerError,
					"server error (record not found)",
					nil, nil))
			}
			return c.JSON(helper.ReponseFormatWithMeta(http.StatusInternalServerError,
				"server error (unexpected)",
				nil, nil))
		}

		var filteredData []cc.CcCore

		if searchInTitle != "" {
			for _, ccData := range ccDatas {
				if strings.Contains(strings.ToLower(ccData.Title), strings.ToLower(searchInTitle)) {
					filteredData = append(filteredData, ccData)
				}
			}
		} else {
			filteredData = ccDatas
		}

		if searchInFrom != "" {
			for _, ccData := range ccDatas {
				if strings.Contains(strings.ToLower(ccData.Title), strings.ToLower(searchInFrom)) {
					filteredData = append(filteredData, ccData)
				}
			}
		}
		if searchInTo != "" {
			for _, ccData := range ccDatas {
				if strings.Contains(strings.ToLower(ccData.Title), strings.ToLower(searchInTo)) {
					filteredData = append(filteredData, ccData)
				}
			}
		}
		if searchInType != "" {
			for _, ccData := range ccDatas {
				if strings.Contains(strings.ToLower(ccData.Title), strings.ToLower(searchInType)) {
					filteredData = append(filteredData, ccData)
				}
			}
		}

		for _, data := range filteredData {
			tmp := Response{
				SubmissionID:   data.SubmisisonID,
				Title:          data.Title,
				SubmissionType: data.SubmissionType,
				Attachment:     data.Attachment,
				From: Sender{
					Name:     data.From.Name,
					Position: data.From.Position,
				},
				To: Receiver{
					Name:     data.To.Name,
					Position: data.To.Position,
				},
			}
			response = append(response, tmp)
		}
		if offsetInt+limitInt > len(response) {
			limitInt = len(filteredData) - offsetInt
		}
		totalData := len(response)
		totalPage := 1
		if len(response) > 0 {
			totalPage = int(math.Ceil(float64(totalData) / float64(limitInt)))
		}
		currentPage := int(math.Ceil(float64(offsetInt+1) / float64(limitInt)))
		if currentPage > totalPage {
			currentPage = totalPage
		}
		meta := Meta{
			CurrentLimit:  limitInt,
			CurrentOffset: offsetInt,
			CurrentPage:   currentPage,
			TotalData:     totalData,
			TotalPage:     totalPage,
		}

		return c.JSON(helper.ReponseFormatWithMeta(http.StatusOK,
			"succes to get ccs data",
			response,
			meta))
	}
}
