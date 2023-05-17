package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ApproveRoutes(e *echo.Echo, pc approve.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/approve", pc.GetSubmissionAprroveHandler(), helper.JWTMiddleWare())
	// e.PUT("/profile", pc.UpdateUserHandler(), helper.JWTMiddleWare())

}
