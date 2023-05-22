package usecase_test

import (
	"errors"
	"testing"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/office/usecase"
	"github.com/stretchr/testify/assert"
)

func TestDeleteOfficeLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	ul := usecase.New(repo)

	t.Run("Succes To Delete Office", func(t *testing.T) {
		repo.On("DeleteOffice", uint(1)).Return(nil).Once()
		err := ul.DeleteOfficeLogic(uint(1))

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error Finding Office", func(t *testing.T) {
		repo.On("DeleteOffice", uint(2)).Return(errors.New("finding office")).Once()
		err := ul.DeleteOfficeLogic(uint(2))

		assert.EqualError(t, err, "bad request, office not found")
		repo.AssertExpectations(t)
	})

	t.Run("Error Deleting Office", func(t *testing.T) {
		repo.On("DeleteOffice", uint(3)).Return(errors.New("cannot delete")).Once()
		err := ul.DeleteOfficeLogic(uint(3))

		assert.EqualError(t, err, "internal server error, cannot delete office")
		repo.AssertExpectations(t)
	})

	t.Run("Other Error", func(t *testing.T) {
		expectedErr := errors.New("some other error")
		repo.On("DeleteOffice", uint(4)).Return(expectedErr).Once()
		err := ul.DeleteOfficeLogic(uint(4))

		assert.EqualError(t, err, expectedErr.Error())
		repo.AssertExpectations(t)
	})
}

func TestGetAllOfficeLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	ul := usecase.New(repo)

	t.Run("Success", func(t *testing.T) {
		expectedOffices := []office.Core{
			{ID: 1, Name: "Office 1"},
			{ID: 2, Name: "Office 2"},
		}

		repo.On("GetAllOffice", 10, 0, "").Return(expectedOffices, nil).Once()
		offices, err := ul.GetAllOfficeLogic(10, 0, "")

		assert.NoError(t, err)
		assert.Equal(t, expectedOffices, offices)
		repo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		expectedErr := errors.New("some error")

		repo.On("GetAllOffice", 10, 0, "").Return([]office.Core{}, expectedErr).Once()
		offices, err := ul.GetAllOfficeLogic(10, 0, "")

		assert.EqualError(t, err, "internal server error")
		assert.Empty(t, offices)
		repo.AssertExpectations(t)
	})
}

func TestAddOfficeLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	ul := usecase.New(repo)

	t.Run("Success", func(t *testing.T) {
		newOffice := office.Core{
			ID:   1,
			Name: "New Office",
		}

		repo.On("InsertOffice", newOffice).Return(nil).Once()
		err := ul.AddOfficeLogic(newOffice)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Column Error", func(t *testing.T) {
		newOffice := office.Core{
			ID:   1,
			Name: "New Office",
		}

		columnErr := errors.New("column error")
		repo.On("InsertOffice", newOffice).Return(columnErr).Once()
		err := ul.AddOfficeLogic(newOffice)

		assert.EqualError(t, err, "server error")
		repo.AssertExpectations(t)
	})

	t.Run("Other Error", func(t *testing.T) {
		newOffice := office.Core{
			ID:   1,
			Name: "New Office",
		}

		otherErr := errors.New("other error")
		repo.On("InsertOffice", newOffice).Return(otherErr).Once()
		err := ul.AddOfficeLogic(newOffice)

		assert.EqualError(t, err, otherErr.Error())
		repo.AssertExpectations(t)
	})
}
