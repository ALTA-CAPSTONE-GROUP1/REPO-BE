package routes

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AdminUserRoutes(e *echo.Echo, uc user.Handler) {
	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.Logger())

	e.POST("/users", uc.RegisterHandler(), helper.JWTMiddleWare())
	e.GET("/users/:id", uc.GetUserByIdHandler(), helper.JWTMiddleWare())
	e.GET("/users", uc.GetAllUserHandler(), helper.JWTMiddleWare())
	e.DELETE("/users/:id", uc.DeleteUserHandler(), helper.JWTMiddleWare())
	e.PUT("/users/:id", uc.UpdateUserHandler(), helper.JWTMiddleWare())
}
