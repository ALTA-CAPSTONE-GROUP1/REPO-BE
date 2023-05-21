package main

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/database"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/routes"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"

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

	profileHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile/handler"
	profileRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile/repository"
	profileLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile/usecase"

	submissionHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/handler"
	submissionRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/repository"
	submissionLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/usecase"

	approveHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve/handler"
	approveRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve/repository"
	approveLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve/usecase"

	ccHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc/handler"
	ccRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc/repository"
	ccLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc/usecase"

	hyperHandler "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve/handler"
	hyperRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve/repository"
	hyperLogic "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve/usecase"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	cfg := config.InitConfig()
	db := database.InitDBMySql(*cfg)
	database.Migrate(db)

	submissioMdl := submissionRepo.New(db)
	submissionSrv := submissionLogic.New(submissioMdl, helper.New())
	submissionCtl := submissionHandler.New(submissionSrv)

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

	profileMdl := profileRepo.New(db)
	profileSrv := profileLogic.New(profileMdl)
	profileCtl := profileHandler.New(profileSrv)

	approveMdl := approveRepo.New(db)
	approveSrv := approveLogic.New(approveMdl)
	approveCtl := approveHandler.New(approveSrv)

	ccMdl := ccRepo.New(db)
	ccSrv := ccLogic.New(ccMdl)
	ccCtl := ccHandler.New(ccSrv)

	hyperMdl := hyperRepo.New(db)
	hyperSrv := hyperLogic.New(hyperMdl)
	hyperCtl := hyperHandler.New(hyperSrv)

	routes.SubTypeRoutes(e, subTypeCtl)
	routes.PositionRoutes(e, pstCtl)
	routes.AdminUserRoutes(e, uAdminCtl)
	routes.AuthRoutes(e, authCtl)
	routes.OfficeRoutes(e, officeCtl)
	routes.ProfileRoutes(e, profileCtl)
	routes.SubmissionRoutes(e, submissionCtl)
	routes.ApproveRoutes(e, approveCtl)
	routes.CcRoutes(e, ccCtl)
	routes.HyperApproveRoutes(e, hyperCtl)

	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal("cannot start server", err.Error())
	}
}
