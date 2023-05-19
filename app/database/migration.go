package database

import (
	adminRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	subRepo "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/repository"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&adminRepo.Users{},
		&adminRepo.Office{},
		&adminRepo.Type{},
		&adminRepo.PositionHasType{},
		&subRepo.Submission{},
		&subRepo.Cc{},
		&subRepo.File{},
		&subRepo.To{},
		&subRepo.Sign{},
	)

	db.SetupJoinTable(&adminRepo.Type{}, "Positions", &adminRepo.PositionHasType{})

	db.Exec(`
		ALTER TABLE position_has_types
		DROP PRIMARY KEY,
		DROP FOREIGN KEY fk_position_has_types_position,
		DROP FOREIGN KEY fk_position_has_types_type;
	`)

	db.Exec(`
	ALTER TABLE position_has_types
	ADD COLUMN id INT PRIMARY KEY AUTO_INCREMENT,
	ADD COLUMN` + "`as`" + ` VARCHAR(10) NOT NULL,
	ADD COLUMN to_level INT,
	ADD COLUMN created_at TIMESTAMP DEFAULT current_timestamp,
	ADD COLUMN deleted_at TIMESTAMP,
	ADD COLUMN value INT,
	ADD CONSTRAINT fk_position_has_types_position FOREIGN KEY (position_id) REFERENCES positions (id) ON DELETE CASCADE,
	ADD CONSTRAINT fk_position_has_types_type FOREIGN KEY (type_id) REFERENCES types (id) ON DELETE CASCADE;
`)

}
