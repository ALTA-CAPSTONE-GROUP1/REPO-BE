package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func HyperApproveRoutes(e *echo.Echo, pc approve.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.PUT("/hyper-approval", pc.UpdateByHyperApprovalHandler(), helper.JWTMiddleWare())
	e.POST("/hyper-approval", pc.GetSubmissionByHyperApprovalHandler(), helper.JWTMiddleWare())

}
