package main

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/database"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/routes"

	uAdminHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user/handler"
	uAdminRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user/repository"
	uAdminLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/user/usecase"

	authHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth/handler"
	authRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth/repository"
	authLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth/usecase"

	stHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype/handler"
	stRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype/repository"
	stLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype/usecase"

	pstHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position/handler"
	pstRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position/repository"
	pstLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position/usecase"

	officeHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office/handler"
	officeRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office/repository"
	officeLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := database.InitDBMySql(*cfg)
	database.Migrate(db)

	subtypeMdl := stRepo.New(db)
	subtypeSrv := stLogic.New(subtypeMdl)
	subTypeCtl := stHandler.New(subtypeSrv)

	pstMdl := pstRepo.New(db)
	pstSrv := pstLogic.New(pstMdl)
	pstCtl := pstHandler.New(pstSrv)

	uAdminMdl := uAdminRepo.New(db)
	uAdminSrv := uAdminLogic.New(uAdminMdl)
	uAdminCtl := uAdminHandler.New(uAdminSrv)

	authMdl := authRepo.New(db)
	authSrv := authLogic.New(authMdl)
	authCtl := authHandler.New(authSrv)

	officeMdl := officeRepo.New(db)
	officeSrv := officeLogic.New(officeMdl)
	officeCtl := officeHandler.New(officeSrv)

	routes.SubTypeRoutes(e, subTypeCtl)
	routes.PositionRoutes(e, pstCtl)
	routes.AdminUserRoutes(e, uAdminCtl)
	routes.AuthRoutes(e, authCtl)
	routes.OfficeRoutes(e, officeCtl)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal("cannot start server", err.Error())
	}
}
