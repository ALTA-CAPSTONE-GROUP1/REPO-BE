package usecase_test

// import (
// 	"errors"
// 	"testing"

// 	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve"
// 	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve/mocks"
// 	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/approve/usecase"
// 	"github.com/stretchr/testify/assert"
// )

// func TestUpdateApprove(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	al := usecase.New(repo)

// 	t.Run("Success", func(t *testing.T) {
// 		userID := "ADMIN"
// 		updateInput := approve.Core{
// 			ID:     1,
// 			Status: "approved",
// 		}

// 		repo.On("UpdateByHyperApproval", userID, updateInput).Return(nil).Once()
// 		err := al.UpdateByHyperApproval(userID, updateInput)

// 		assert.NoError(t, err)
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Invalid Password", func(t *testing.T) {
// 		userID := "ADMIN"
// 		updateInput := approve.Core{
// 			ID:     2,
// 			Status: "rejected",
// 		}

// 		errMsg := "hashing password error"
// 		repo.On("UpdateByHyperApproval", userID, updateInput).Return(errors.New(errMsg)).Once()
// 		err := al.UpdateByHyperApproval(userID, updateInput)

// 		assert.EqualError(t, err, "is invalid")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("No Rows Affected", func(t *testing.T) {
// 		userID := "ADMIN"
// 		updateInput := approve.Core{
// 			ID:     3,
// 			Status: "revised",
// 		}

// 		errMsg := "no rows affected on update submission"
// 		repo.On("UpdateByHyperApproval", userID, updateInput).Return(errors.New(errMsg)).Once()
// 		err := al.UpdateByHyperApproval(userID, updateInput)

// 		assert.EqualError(t, err, "data is up to date")
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Other Error", func(t *testing.T) {
// 		userID := "ADMIN"
// 		updateInput := approve.Core{
// 			ID:     1,
// 			Status: "approved",
// 		}

// 		errMsg := "other error"
// 		repo.On("UpdateByHyperApproval", userID, updateInput).Return(errors.New(errMsg)).Once()
// 		err := al.UpdateByHyperApproval(userID, updateInput)

// 		assert.EqualError(t, err, errMsg)
// 		repo.AssertExpectations(t)
// 	})
// }

// func TestGetSubmissionById(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	al := usecase.New(repo)

// 	t.Run("Success", func(t *testing.T) {
// 		userID := "ADMIN"
// 		id := 1
// 		token := "admin123"
// 		expectedSubmission := approve.GetSubmissionByIDCore{
// 			Status: "pending",
// 		}

// 		repo.On("SelectSubmissionByHyperApproval", userID, id, token).Return(expectedSubmission, nil).Once()
// 		submission, err := al.GetSubmissionByHyperApproval(userID, id, token)

// 		assert.NoError(t, err)
// 		assert.Equal(t, expectedSubmission, submission)
// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Error", func(t *testing.T) {
// 		userID := "ADMIN"
// 		id := 1
// 		token := "admin123"
// 		errMsg := "failed to find submission error"

// 		repo.On("SelectSubmissionByHyperApproval", userID, id, token).Return(approve.GetSubmissionByIDCore{}, errors.New(errMsg)).Once()
// 		submission, err := al.GetSubmissionByHyperApproval(userID, id, token)

// 		assert.EqualError(t, err, "internal server error")
// 		assert.Equal(t, approve.GetSubmissionByIDCore{}, submission)
// 		repo.AssertExpectations(t)
// 	})
// }
