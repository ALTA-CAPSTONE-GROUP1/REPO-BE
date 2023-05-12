package usecase_test

import (
	"errors"
	"testing"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin/position/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAddPosition(t *testing.T) {
	repo := mocks.NewRepository(t)
	ul := usecase.New(repo)
	rightPosition := position.Core{
		Name: "Purchasing Manager",
		Tag:  "PURM",
	}
	t.Run("Succes Create Position", func(t *testing.T) {
		repo.On("InsertPosition", rightPosition).Return(nil)
		err := ul.AddPositionLogic(rightPosition)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Create Position - Server Error", func(t *testing.T) {
		errPosition := position.Core{}
		repo.On("InsertPosition", errPosition).Return(errors.New("column"))
		err := ul.AddPositionLogic(errPosition)

		assert.Error(t, err)
		assert.EqualError(t, err, "server error")
		repo.AssertExpectations(t)
	})

	t.Run("Failed Create Position - Unexpected Error", func(t *testing.T) {
		errPosition := position.Core{
			Name: "dsiandia",
			Tag:  "SJBAU",
		}
		repo.On("InsertPosition", errPosition).Return(errors.New("too many arguments for record"))
		err := ul.AddPositionLogic(errPosition)

		assert.Error(t, err)
		assert.EqualError(t, err, "too many arguments for record")
		repo.AssertExpectations(t)
	})
}
