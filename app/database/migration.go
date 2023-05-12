package database

import (
	adminRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	subRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/repository"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(adminRepo.Users{})
	db.AutoMigrate(adminRepo.Office{})
	db.AutoMigrate(adminRepo.Position{})
	db.AutoMigrate(adminRepo.PositionHasType{})
	db.AutoMigrate(subRepo.Submission{})
	db.AutoMigrate(subRepo.Cc{})
	db.AutoMigrate(subRepo.File{})
	db.AutoMigrate(subRepo.To{})
	db.AutoMigrate(subRepo.Sign{})
}
