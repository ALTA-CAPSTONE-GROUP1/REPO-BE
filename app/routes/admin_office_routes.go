package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func OfficeRoutes(e *echo.Echo, pc office.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/office", pc.AddOfficeHandler(), helper.JWTMiddleWare())
	e.GET("/office", pc.GetAllOfficeHandler(), helper.JWTMiddleWare())
	e.DELETE("/office/:id", pc.DeleteOfficeHandler(), helper.JWTMiddleWare())

}
