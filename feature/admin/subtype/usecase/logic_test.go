package usecase_test

import (
	"errors"
	"fmt"
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

func TestGetSubTypesLogic(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	sl := usecase.New(mockRepo)

	// Mock data
	limit := 10
	offset := 0
	search := "test"
	subtypeData := []subtype.GetSubmissionTypeCore{
		{
			SubmissionTypeName: "Test1",
			Value:              5000000,
			Requirement:        "Test Requirement 1",
		},
		{
			SubmissionTypeName: "Test2",
			Value:              7500000,
			Requirement:        "Test Requirement 2",
		},
	}
	positionData := []subtype.GetPosition{
		{
			PositionName: "Owner1",
			PositionTag:  "owner-1",
		},
		{
			PositionName: "Owner2",
			PositionTag:  "owner-2",
		},
	}

	t.Run("Success", func(t *testing.T) {
		// Mock repository function calls
		mockRepo.On("GetSubTypes", limit, offset, search).Return(subtypeData, positionData, nil)

		// Call the function
		resultSubtypeData, resultPositionData, err := sl.GetSubTypesLogic(limit, offset, search)

		// Verify the results
		assert.NoError(t, err)
		assert.Equal(t, subtypeData, resultSubtypeData)
		assert.Equal(t, positionData, resultPositionData)

		// Verify that the mock repository functions were called
		mockRepo.AssertExpectations(t)
	})

	t.Run("NegativeLimit", func(t *testing.T) {
		// Call the function with negative limit
		resultSubtypeData, resultPositionData, err := sl.GetSubTypesLogic(-1, offset, search)

		// Verify the error
		assert.EqualError(t, err, "cannot accept limit value = -1")
		assert.Nil(t, resultSubtypeData)
		assert.Nil(t, resultPositionData)

		// Verify that the mock repository functions were not called
		mockRepo.AssertExpectations(t)
	})

	t.Run("FailedToRetrievePositions", func(t *testing.T) {
		// Mock repository function calls
		mockRepo.On("GetSubTypes", 2, 3, "test2").Return(nil, nil, fmt.Errorf("finding all positions")).Once()

		// Call the function
		resultSubtypeData, resultPositionData, err := sl.GetSubTypesLogic(2, 3, "test2")

		// Verify the error
		assert.ErrorContains(t, err, "failed to retrieve positions. finding all positions")
		assert.Nil(t, resultSubtypeData)
		assert.Nil(t, resultPositionData)

		// Verify that the mock repository functions were called
		mockRepo.AssertExpectations(t)
	})

	t.Run("FailedToRetrieveSubmissionTypes", func(t *testing.T) {
		// Mock repository function calls
		mockRepo.On("GetSubTypes", 1, 10, "search3").Return(nil, nil, fmt.Errorf("all submission types")).Once()

		// Call the function
		resultSubtypeData, resultPositionData, err := sl.GetSubTypesLogic(1, 10, "search3")

		// Verify the error
		assert.EqualError(t, err, "failed to retrieve submission types. all submission types")
		assert.Nil(t, resultSubtypeData)
		assert.Nil(t, resultPositionData)

		// Verify that the mock repository functions were called
		mockRepo.AssertExpectations(t)
	})

	t.Run("FailedToRetrievePositionHasTypes", func(t *testing.T) {
		// Mock repository function calls
		mockRepo.On("GetSubTypes", 5, 10, "search4").Return(nil, nil, fmt.Errorf("all position_has_types")).Once()

		// Call the function
		resultSubtypeData, resultPositionData, err := sl.GetSubTypesLogic(5, 10, "search4")

		// Verify the error
		assert.EqualError(t, err, "failed to retrieve position_has_types. all position_has_types")
		assert.Nil(t, resultSubtypeData)
		assert.Nil(t, resultPositionData)

		// Verify that the mock repository functions were called
		mockRepo.AssertExpectations(t)
	})

	t.Run("FailedToRetrievePositionHasTypes", func(t *testing.T) {
		// Mock repository function calls
		mockRepo.On("GetSubTypes", 2, 3, "searchlain").Return(nil, nil, fmt.Errorf("unexpected whatever it is")).Once()

		// Call the function
		resultSubtypeData, resultPositionData, err := sl.GetSubTypesLogic(2, 3, "searchlain")

		// Verify the error
		assert.EqualError(t, err, "failed to get submission types with unexpected error. unexpected whatever it is")
		assert.Nil(t, resultSubtypeData)
		assert.Nil(t, resultPositionData)

		// Verify that the mock repository functions were called
		mockRepo.AssertExpectations(t)
	})
}
