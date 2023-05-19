package repository

import (
	"fmt"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/labstack/gommon/log"
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

func (st *subTypeModel) InsertSubType(req subtype.RepoData) error {
	typeData := admin.Type{
		Name:        req.TypeName,
		Requirement: req.TypeRequirement,
	}

	tx := st.db.Begin()

	if err := tx.Create(&typeData).Error; err != nil {
		tx.Rollback()
		log.Errorf("failed to insert type data: %v", err)
		return err
	}

	owners, err := st.getPositionByTag(tx, req.OwnersTag)
	if err != nil {
		tx.Rollback()
		log.Errorf("owners position not found = %v", req.OwnersTag)
		return err
	}

	values := make([]admin.PositionHasType, 0, len(req.SubTypeInterdependence)*(len(owners)+len(req.SubTypeInterdependence[0].TosTag)+len(req.SubTypeInterdependence[0].CcsTag)))

	for _, SubmissionValue := range req.SubTypeInterdependence {
		value := SubmissionValue.Value

		approverPositions, err := st.getPositionByTag(tx, SubmissionValue.TosTag)
		if err != nil {
			tx.Rollback()
			log.Error("cannot find authorized officials approver by tag")
			return err
		}

		ccPositions, err := st.getPositionByTag(tx, SubmissionValue.CcsTag)
		if err != nil {
			tx.Rollback()
			log.Error("cannot find authorized officials ccs by tag")
			return err
		}

		for _, owner := range owners {
			newOwner := admin.PositionHasType{
				TypeID:     typeData.ID,
				PositionID: owner.ID,
				Value:      value,
				As:         "Owner",
			}

			values = append(values, newOwner)
		}
		level := 1
		for _, approver := range approverPositions {
			newApprover := admin.PositionHasType{
				TypeID:     typeData.ID,
				PositionID: approver.ID,
				Value:      value,
				As:         "To",
				ToLevel:    level,
			}

			values = append(values, newApprover)

			level++
		}

		for _, cc := range ccPositions {
			newCc := admin.PositionHasType{
				TypeID:     typeData.ID,
				PositionID: cc.ID,
				Value:      value,
				As:         "CC",
			}

			values = append(values, newCc)
		}
	}

	if len(values) > 0 {
		if err := tx.Create(&values).Error; err != nil {
			tx.Rollback()
			log.Errorf("failed to insert position has type data: %v", err)
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Errorf("failed to commit transaction: %v", err)
		return err
	}

	return nil
}
func (st *subTypeModel) getPositionByTag(tx *gorm.DB, positionTags []string) ([]admin.Position, error) {
	var positionData []admin.Position

	for _, val := range positionTags {
		var position admin.Position
		if err := tx.Where("tag = ?", val).First(&position).Error; err != nil {
			log.Errorf("failed to get position data: %v", err)
			return []admin.Position{}, fmt.Errorf("cannot find this position %s", val)
		}

		positionData = append(positionData, position)
	}

	return positionData, nil
}

func (st *subTypeModel) GetSubTypes(limit int, offset int, search string) ([]subtype.GetSubmissionTypeCore, []subtype.GetPosition, error) {
	var dbpositions []admin.Position
	var resPositions []subtype.GetPosition
	var submissionTypeCoreData []subtype.GetSubmissionTypeCore
	var types []admin.Type
	var hasTypes []admin.PositionHasType

	if err := st.db.Find(&dbpositions).Error; err != nil {
		log.Errorf("failed on finding all positions %w", err)
		return nil, nil, fmt.Errorf("error on finding all positions in get all submission type %w", err)
	}

	for _, dbposition := range dbpositions {
		tmpPosition := subtype.GetPosition{
			PositionName: dbposition.Name,
			PositionTag:  dbposition.Tag,
		}
		resPositions = append(resPositions, tmpPosition)
	}

	if err := st.db.Find(&types).Error; err != nil {
		log.Errorf("failed on finding all submission types %w", err)
		return nil, nil, fmt.Errorf("error when get all submission type %w", err)
	}

	if err := st.db.Find(&hasTypes).Error; err != nil {
		log.Errorf("failed on finding all position_has_types for getall submission types %w", err)
		return nil, nil, fmt.Errorf("error when get all positions_has_types %w", err)
	}

	for _, t := range types {
		for _, h := range hasTypes {
			if t.ID == h.TypeID && h.As == "Owner" {
				tmp := subtype.GetSubmissionTypeCore{
					SubmissionTypeName: t.Name,
					Value:              h.Value,
					Requirement:        t.Requirement,
				}
				submissionTypeCoreData = append(submissionTypeCoreData, tmp)
			}
		}
	}
	fmt.Println(submissionTypeCoreData)
	return submissionTypeCoreData, resPositions, nil
}

func (st *subTypeModel) DeleteSubType(subTypeName string) error {
	if err := st.db.Where("name = ?", subTypeName).Delete(&admin.Type{}).Error; err != nil {
		log.Errorf("error on delete subtype by name, %w", err)
		return fmt.Errorf("failed to delete subtype by name %w", err)
	}
	return nil
}
