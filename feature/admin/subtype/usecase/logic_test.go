package usecase_test

import (
	"errors"
	"testing"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/subtype/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAddSubTypeLogic(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	sl := usecase.New(mockRepo)
	succesData := subtype.RepoData{
		TypeName:        "Test",
		TypeRequirement: "Test Requirement",
		OwnersTag:       []string{"owner-1"},
		SubTypeInterdependence: []subtype.RepoDataInterdependence{
			{
				Value:  5000000,
				TosTag: []string{"to-1"},
				CcsTag: []string{"cc-1"},
			},
		},
	}
	t.Run("Succes Create Position", func(t *testing.T) {
		mockRepo.On("InsertSubType", succesData).Return(nil).Once()

		err := sl.AddSubTypeLogic(subtype.Core{
			SubmissionTypeName: "Test",
			Requirement:        "Test Requirement",
			PositionTag:        []string{"owner-1"},
			SubmissionValues: []subtype.ValueDetails{
				{
					Value:         5000000,
					TagPositionTo: []string{"to-1"},
					TagPositionCC: []string{"cc-1"},
				},
			},
		})

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Succes Create Position", func(t *testing.T) {
		mockRepo.On("InsertSubType", succesData).Return(errors.New("failed to insert type data")).Once()

		err := sl.AddSubTypeLogic(subtype.Core{
			SubmissionTypeName: "Test",
			Requirement:        "Test Requirement",
			PositionTag:        []string{"owner-1"},
			SubmissionValues: []subtype.ValueDetails{
				{
					Value:         5000000,
					TagPositionTo: []string{"to-1"},
					TagPositionCC: []string{"cc-1"},
				},
			},
		})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to insert submission type data")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Succes Create Position", func(t *testing.T) {
		mockRepo.On("InsertSubType", succesData).Return(errors.New("owners position not found")).Once()

		err := sl.AddSubTypeLogic(subtype.Core{
			SubmissionTypeName: "Test",
			Requirement:        "Test Requirement",
			PositionTag:        []string{"owner-1"},
			SubmissionValues: []subtype.ValueDetails{
				{
					Value:         5000000,
					TagPositionTo: []string{"to-1"},
					TagPositionCC: []string{"cc-1"},
				},
			},
		})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to add user as authorized to make")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Succes Create Position", func(t *testing.T) {
		mockRepo.On("InsertSubType", succesData).Return(errors.New("cannot find authorized officials approver by tag")).Once()

		err := sl.AddSubTypeLogic(subtype.Core{
			SubmissionTypeName: "Test",
			Requirement:        "Test Requirement",
			PositionTag:        []string{"owner-1"},
			SubmissionValues: []subtype.ValueDetails{
				{
					Value:         5000000,
					TagPositionTo: []string{"to-1"},
					TagPositionCC: []string{"cc-1"},
				},
			},
		})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to add approver to the database")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Succes Create Position", func(t *testing.T) {
		mockRepo.On("InsertSubType", succesData).Return(errors.New("failed to insert position has type data")).Once()

		err := sl.AddSubTypeLogic(subtype.Core{
			SubmissionTypeName: "Test",
			Requirement:        "Test Requirement",
			PositionTag:        []string{"owner-1"},
			SubmissionValues: []subtype.ValueDetails{
				{
					Value:         5000000,
					TagPositionTo: []string{"to-1"},
					TagPositionCC: []string{"cc-1"},
				},
			},
		})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to add roles to data type")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Succes Create Position", func(t *testing.T) {
		mockRepo.On("InsertSubType", succesData).Return(errors.New("failed to commit transaction")).Once()

		err := sl.AddSubTypeLogic(subtype.Core{
			SubmissionTypeName: "Test",
			Requirement:        "Test Requirement",
			PositionTag:        []string{"owner-1"},
			SubmissionValues: []subtype.ValueDetails{
				{
					Value:         5000000,
					TagPositionTo: []string{"to-1"},
					TagPositionCC: []string{"cc-1"},
				},
			},
		})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to save all data to database (commit error)")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Succes Create Position", func(t *testing.T) {
		mockRepo.On("InsertSubType", succesData).Return(errors.New("some error that doesnt used 4545487887")).Once()

		err := sl.AddSubTypeLogic(subtype.Core{
			SubmissionTypeName: "Test",
			Requirement:        "Test Requirement",
			PositionTag:        []string{"owner-1"},
			SubmissionValues: []subtype.ValueDetails{
				{
					Value:         5000000,
					TagPositionTo: []string{"to-1"},
					TagPositionCC: []string{"cc-1"},
				},
			},
		})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "unexpected error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Succes Create Position", func(t *testing.T) {
		mockRepo.On("InsertSubType", succesData).Return(errors.New("cannot find authorized officials ccs by tag")).Once()

		err := sl.AddSubTypeLogic(subtype.Core{
			SubmissionTypeName: "Test",
			Requirement:        "Test Requirement",
			PositionTag:        []string{"owner-1"},
			SubmissionValues: []subtype.ValueDetails{
				{
					Value:         5000000,
					TagPositionTo: []string{"to-1"},
					TagPositionCC: []string{"cc-1"},
				},
			},
		})

		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "failed to add cc to the database")
		mockRepo.AssertExpectations(t)
	})
}
