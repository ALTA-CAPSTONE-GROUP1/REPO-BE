package usecase_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/auth/usecase"
	"github.com/stretchr/testify/assert"
)

func TestSignValidationLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	usecaselogic := usecase.New(repo)

	t.Run("Record Not Found", func(t *testing.T) {
		notFoundID := "abas123"
		repo.On("SignVaidation", notFoundID).Return(auth.SignCore{}, errors.New("sign record")).Once()
		result, err := usecaselogic.SignVallidationLogic(notFoundID)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.EqualError(t, err, "sign record not found")
		repo.AssertExpectations(t)
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		notFoundID := "abas1223"
		repo.On("SignVaidation", notFoundID).Return(auth.SignCore{}, errors.New("anything")).Once()
		result, err := usecaselogic.SignVallidationLogic(notFoundID)

		assert.Error(t, err)
		assert.EqualError(t, err, "server error, unexpected")
		assert.Empty(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Succes To Get Data", func(t *testing.T) {
		notFoundID := "abas1234"
		dbresSucces := auth.SignCore{
			Title:            "apa saja",
			Officialname:     "Kholil",
			Officialposition: "Manager",
			Date:             time.Now().String(),
		}
		repo.On("SignVaidation", notFoundID).Return(dbresSucces, nil).Once()
		result, err := usecaselogic.SignVallidationLogic(notFoundID)

		assert.NoError(t, err)
		assert.NotEmpty(t, result, dbresSucces)
		repo.AssertExpectations(t)
	})
}

func TestLogInLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	ul := usecase.New(repo)

	t.Run("IDCannotBeBlank", func(t *testing.T) {
		id := ""
		password := "kadal12345"
		repo.On("Login", id, password).Return(auth.Core{}, errors.New("user does not exist"))

		_, err := ul.LogInLogic(id, password)
		assert.Error(t, err)
		assert.EqualError(t, err, "id cannot be blank")

		repo.AssertExpectations(t)
	})

	t.Run("PasswordIsWrong", func(t *testing.T) {
		id := "AM2"
		password := "passwordsalah"
		repo.On("Login", id, password).Return(auth.Core{}, errors.New("password is wrong"))

		_, err := ul.LogInLogic(id, password)
		assert.Error(t, err)
		assert.EqualError(t, err, "password is wrong")

		repo.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		id := "AM3"
		password := "passwordsalahbanget"
		repo.On("Login", id, password).Return(auth.Core{}, errors.New("internal server error"))

		_, err := ul.LogInLogic(id, password)
		assert.Error(t, err)
		assert.EqualError(t, err, "internal server error")

		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		id := "AM5"
		password := "passwordBenar"
		expectedCore := auth.Core{
			ID:          "AM5",
			OfficeID:    1,
			PositionID:  2,
			Name:        "Kholil",
			Email:       "kholil@gmail.com",
			PhoneNumber: "081223536464",
			Password:    "passwordbenar",
		}
		repo.On("Login", id, password).Return(expectedCore, nil)

		resultNoErrror, err := ul.LogInLogic(id, password)
		assert.NoError(t, err)
		assert.Equal(t, expectedCore, resultNoErrror)

		repo.AssertExpectations(t)
	})
}
