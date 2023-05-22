package usecase_test

import (
	"errors"
	"testing"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/approve/usecase"
	"github.com/stretchr/testify/assert"
)

func TestUpdateApprove(t *testing.T) {
	repo := mocks.NewRepository(t)
	al := usecase.New(repo)

	t.Run("Success", func(t *testing.T) {
		userID := "AM1"
		id := 1
		updateInput := approve.Core{
			Status: "approved",
		}

		repo.On("UpdateApprove", userID, id, updateInput).Return(nil).Once()
		err := al.UpdateApprove(userID, id, updateInput)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		userID := "AM2"
		id := 1
		updateInput := approve.Core{
			Status: "approved",
		}

		errMsg := "hashing password error"
		repo.On("UpdateApprove", userID, id, updateInput).Return(errors.New(errMsg)).Once()
		err := al.UpdateApprove(userID, id, updateInput)

		assert.EqualError(t, err, "is invalid")
		repo.AssertExpectations(t)
	})

	t.Run("No Rows Affected", func(t *testing.T) {
		userID := "AM3"
		id := 1
		updateInput := approve.Core{
			Status: "approved",
		}

		errMsg := "no rows affected on update submission"
		repo.On("UpdateApprove", userID, id, updateInput).Return(errors.New(errMsg)).Once()
		err := al.UpdateApprove(userID, id, updateInput)

		assert.EqualError(t, err, "data is up to date")
		repo.AssertExpectations(t)
	})

	t.Run("Other Error", func(t *testing.T) {
		userID := "AM4"
		id := 1
		updateInput := approve.Core{
			Status: "approved",
		}

		errMsg := "other error"
		repo.On("UpdateApprove", userID, id, updateInput).Return(errors.New(errMsg)).Once()
		err := al.UpdateApprove(userID, id, updateInput)

		assert.EqualError(t, err, errMsg)
		repo.AssertExpectations(t)
	})
}

func TestGetSubmissionById(t *testing.T) {
	repo := mocks.NewRepository(t)
	al := usecase.New(repo)

	t.Run("Success", func(t *testing.T) {
		userID := "rm1"
		id := 1
		expectedSubmission := approve.Core{
			ID:     id,
			UserID: userID,
			Status: "pending",
		}

		repo.On("SelectSubmissionById", userID, id).Return(expectedSubmission, nil).Once()
		submission, err := al.GetSubmissionById(userID, id)

		assert.NoError(t, err)
		assert.Equal(t, expectedSubmission, submission)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		userID := "rm2"
		id := 1
		errMsg := "failed to find submission error"

		repo.On("SelectSubmissionById", userID, id).Return(approve.Core{}, errors.New(errMsg)).Once()
		submission, err := al.GetSubmissionById(userID, id)

		assert.EqualError(t, err, "internal server error")
		assert.Equal(t, approve.Core{}, submission)
		repo.AssertExpectations(t)
	})
}

func TestGetSubmissionApprove(t *testing.T) {
	repo := mocks.NewRepository(t)
	al := usecase.New(repo)

	t.Run("Success", func(t *testing.T) {
		userID := "RM22"
		searchParams := approve.GetAllQueryParams{
			Limit:  10,
			Offset: 0,
		}
		expectedSubmissions := []approve.Core{
			{
				ID:     1,
				UserID: userID,
			},
			{
				ID:     2,
				UserID: userID,
			},
		}

		repo.On("SelectSubmissionAprrove", userID, searchParams).Return(expectedSubmissions, nil).Once()
		result, err := al.GetSubmissionAprrove(userID, searchParams)

		assert.NoError(t, err)
		assert.Equal(t, expectedSubmissions, result)
		repo.AssertExpectations(t)
	})

	t.Run("Error Internal Server", func(t *testing.T) {
		userID := "RM22"
		searchParams := approve.GetAllQueryParams{
			Limit:  10,
			Offset: 0,
		}
		expectedError := []approve.Core{}

		repo.On("SelectSubmissionAprrove", userID, searchParams).Return([]approve.Core{}, errors.New("anything error")).Once()

		result, err := al.GetSubmissionAprrove(userID, searchParams)

		assert.Error(t, err)
		assert.Equal(t, expectedError, result)
		repo.AssertExpectations(t)
	})

}
