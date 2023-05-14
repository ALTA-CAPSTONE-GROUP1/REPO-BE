package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SubTypeRoutes(e *echo.Echo, stc subtype.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/submission-type", stc.AddTypeHandler(), helper.JWTMiddleWare())
	e.GET("/submission-type", stc.GetTypesHandler(), helper.JWTMiddleWare())
	e.DELETE("/submission-type", stc.DeleteTypeHandler(), helper.JWTMiddleWare())

}
