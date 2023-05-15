package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ProfileRoutes(e *echo.Echo, pc profile.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.GET("/profile", pc.ProfileHandler(), helper.JWTMiddleWare())
	e.PUT("/profile", pc.UpdateUserHandler(), helper.JWTMiddleWare())

}
