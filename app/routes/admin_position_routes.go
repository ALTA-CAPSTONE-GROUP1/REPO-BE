package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func PositionRoutes(e *echo.Echo, pc position.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/submission-type", pc.AddPositionHandler(), helper.JWTMiddleWare())
	e.GET("/submission-type", pc.GetAllPositionHandler(), helper.JWTMiddleWare())
	e.DELETE("/submission-type", pc.DeletePositionHandler(), helper.JWTMiddleWare())

}
