package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SubmissionRoutes(e *echo.Echo, sc submission.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/submission", sc.AddSubmissionHandler(), helper.JWTMiddleWare())
	e.GET("/submission", sc.GetAllSubmissionHandler(), helper.JWTMiddleWare())
	e.GET("submission/requirements", sc.FindRequirementHandler(), helper.JWTMiddleWare())
	e.GET("/submissions/:submission_id", sc.GetSubmissionByIdHandler(), helper.JWTMiddleWare())
}
