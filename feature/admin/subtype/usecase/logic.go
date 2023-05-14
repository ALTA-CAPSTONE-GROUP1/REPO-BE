package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/labstack/gommon/log"
)

type subTypeLogic struct {
	st subtype.Repository
}

func New(st subtype.Repository) subtype.UseCase {
	return &subTypeLogic{
		st: st,
	}
}

func (stl *subTypeLogic) AddSubTypeLogic(newType subtype.Core) error {
	var insertData subtype.RepoData

	insertData.TypeName = newType.SubmissionTypeName
	insertData.TypeRequirement = newType.Requirement

	insertData.OwnersTag = append(insertData.OwnersTag, newType.PositionTag...)

	for _, submissionValue := range newType.SubmissionValues {
		var interdependence subtype.RepoDataInterdependence
		interdependence.Value = submissionValue.Value

		interdependence.TosTag = append(interdependence.TosTag, submissionValue.TagPositionTo...)
		interdependence.CcsTag = append(interdependence.CcsTag, submissionValue.TagPositionCC...)

		insertData.SubTypeInterdependence = append(insertData.SubTypeInterdependence, interdependence)
	}

	if err := stl.st.InsertSubType(insertData); err != nil {
		if strings.Contains(err.Error(), "failed to insert type data") {
			return fmt.Errorf("failed to insert submission type data %w", err)
		} else if strings.Contains(err.Error(), "owners position not found") {
			return fmt.Errorf("failed to add user as authorized to make this submission type %w", err)
		} else if strings.Contains(err.Error(), "cannot find authorized officials approver by tag") {
			return fmt.Errorf("failed to add approver to the database %w", err)
		} else if strings.Contains(err.Error(), "cannot find authorized officials ccs by tag") {
			return fmt.Errorf("failed to add cc to the database %w", err)
		} else if strings.Contains(err.Error(), "failed to insert position has type data") {
			return fmt.Errorf("failed to add roles to data type %w", err)
		} else if strings.Contains(err.Error(), "failed to commit transaction") {
			return fmt.Errorf("failed to save all data to database (commit error) %w", err)
		} else {
			return fmt.Errorf("unexpected error %w", err)
		}
	}

	return nil
}

func (stl *subTypeLogic) GetSubTypesLogic(limit int, offset int, search string) ([]subtype.GetSubmissionTypeCore, []subtype.GetPosition, error) {
	if limit <= 0 {
		log.Errorf("cannot accept negative value of limit query")
		return nil, nil, fmt.Errorf("cannot accept limit value = %d", limit)
	}

	subtypeData, positionData, err := stl.st.GetSubTypes(limit, offset, search)
	if err != nil {
		if strings.Contains(err.Error(), "finding all positions") {
			log.Errorf("Failed to retrieve positions: %v", err)
			return nil, nil, fmt.Errorf("failed to retrieve positions. %v", err)
		}
		if strings.Contains(err.Error(), "all submission types") {
			log.Errorf("Failed to retrieve submission types: %v", err)
			return nil, nil, fmt.Errorf("failed to retrieve submission types. %v", err)
		}
		if strings.Contains(err.Error(), "all position_has_types") {
			log.Errorf("Failed to retrieve position_shas_types: %v", err)
			return nil, nil, fmt.Errorf("failed to retrieve position_has_types. %v", err)
		}

		log.Errorf("Failed to get submission types (unexpected error): %v", err)
		return nil, nil, fmt.Errorf("failed to get submission types with unexpected error. %v", err)
	}

	return subtypeData, positionData, nil
}

func (stl *subTypeLogic) DeleteSubTypeLogic(subTypeName string) error {
	if err := stl.st.DeleteSubType(subTypeName); err != nil {
		log.Error("error on calling delete SubType in logic")
		if strings.Contains(err.Error(), "empty_set") {
			return errors.New("subtypename not found")
		} else if strings.Contains(err.Error(), "failed to find subtypename for delete") {
			return errors.New("subtypename not found")
		} else if strings.Contains(err.Error(), "failed to delete subtype by name") {
			return errors.New("error on delete the subtype by name")
		} else {
			log.Error("unexpected error")
			return fmt.Errorf("unexpected error, %w", err)
		}
	}
	return nil
}
