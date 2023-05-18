package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CcRoutes(e *echo.Echo, ch cc.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/cc", ch.GetAllCcHander(), helper.JWTMiddleWare())

}
