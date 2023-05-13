package usecase

import (
	"fmt"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
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
