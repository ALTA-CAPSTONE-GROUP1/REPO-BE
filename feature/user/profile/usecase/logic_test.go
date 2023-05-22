package usecase_test

import (
	"errors"
	"testing"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/profile/usecase"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	repo := mocks.NewRepository(t)
	pl := usecase.New(repo)

	t.Run("Success", func(t *testing.T) {
		userID := "12345"
		updateUser := profile.Core{
			Name:  "adi yuda",
			Email: "adiyuda@gmail.com",
		}

		repo.On("UpdateUser", userID, updateUser).Return(nil).Once()
		err := pl.UpdateUser(userID, updateUser)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Invalid Password", func(t *testing.T) {
		userID := "12345"
		updateUser := profile.Core{
			Name:  "bono",
			Email: "boono@gmail.com",
		}

		errMsg := "hashing password error"
		repo.On("UpdateUser", userID, updateUser).Return(errors.New(errMsg)).Once()
		err := pl.UpdateUser(userID, updateUser)

		assert.EqualError(t, err, "is invalid")
		repo.AssertExpectations(t)
	})

	t.Run("No Rows Affected", func(t *testing.T) {
		userID := "12345"
		updateUser := profile.Core{
			Name: "bana",
		}

		errMsg := "no rows affected on update user"
		repo.On("UpdateUser", userID, updateUser).Return(errors.New(errMsg)).Once()
		err := pl.UpdateUser(userID, updateUser)

		assert.EqualError(t, err, "data is up to date")
		repo.AssertExpectations(t)
	})

	t.Run("Other Error", func(t *testing.T) {
		userID := "12345"
		updateUser := profile.Core{
			Name: "bini",
		}

		errMsg := "other error"
		repo.On("UpdateUser", userID, updateUser).Return(errors.New(errMsg)).Once()
		err := pl.UpdateUser(userID, updateUser)

		assert.EqualError(t, err, errMsg)
		repo.AssertExpectations(t)
	})
}

func TestProfileLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	pl := usecase.New(repo)
	t.Run("Success", func(t *testing.T) {
		userID := "pm12"
		expectedProfile := profile.Core{
			ID:   userID,
			Name: "John Doe",
		}

		repo.On("Profile", userID).Return(expectedProfile, nil).Once()
		profile, err := pl.ProfileLogic(userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedProfile, profile)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		userID := "pm21"
		errMsg := "failed to find user error"

		repo.On("Profile", userID).Return(profile.Core{}, errors.New(errMsg)).Once()
		errProfile, err := pl.ProfileLogic(userID)

		assert.EqualError(t, err, "internal server error")
		assert.Equal(t, profile.Core{}, errProfile)
		repo.AssertExpectations(t)
	})
}
