package repository

import (
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"gorm.io/gorm"
)

type subTypeModel struct {
	db *gorm.DB
}

func New(db *gorm.DB) subtype.Repository {
	return &subTypeModel{
		db: db,
	}
}

func (sbm *subTypeModel) InsertSubType(newType subtype.Core) error {
	
}
