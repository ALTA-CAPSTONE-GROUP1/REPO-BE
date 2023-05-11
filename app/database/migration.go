package database

import (
	aRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/repository"
	apRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/approver/repository"
	auRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth/repository"
	sRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/repository"
	uRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/repository"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(uRepo.User{})
	db.AutoMigrate(sRepo.Submission{})
	db.AutoMigrate(auRepo.Auth{})
	db.AutoMigrate(aRepo.Admin{})
	db.AutoMigrate(apRepo.Approver{})

}
