package usecase_test

import (
	"errors"
	"testing"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/user/cc/usecase"
	"github.com/stretchr/testify/assert"
)

func TestAddSubmissionLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	ccl := usecase.New(repo)

	t.Run("Record Not Found", func(t *testing.T) {
		recordNotFoundID := "RM1"

		repo.On("GetAllCc", recordNotFoundID).Return([]cc.CcCore{}, errors.New("record")).Once()

		result, err := ccl.GetAllCcLogic(recordNotFoundID)

		assert.Empty(t, result)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")

		repo.AssertExpectations(t)
	})

	t.Run("Record Not Found", func(t *testing.T) {
		otherErrorID := "RM2"

		repo.On("GetAllCc", otherErrorID).Return([]cc.CcCore{}, errors.New("anything"))

		result, err := ccl.GetAllCcLogic(otherErrorID)

		assert.Empty(t, result)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "anything")

		repo.AssertExpectations(t)
	})

	t.Run("Succes Get Ccs", func(t *testing.T) {
		succesID := "RM4"
		repoReturn := []cc.CcCore{
			{
				SubmisisonID: 1,
				From: cc.Sender{
					Name:     "Bohang",
					Position: "Security",
				},
				To: cc.Receiver{
					Name:     "Bona",
					Position: "Head Security",
				},
				Title:          "Pengajuan Penambahan Karyawan",
				SubmissionType: "Man Power",
				Attachment:     "facebook.com",
			},
		}

		ccResult := []cc.CcCore{
			{
				SubmisisonID: 1,
				From: cc.Sender{
					Name:     "Bohang",
					Position: "Security",
				},
				To: cc.Receiver{
					Name:     "Bona",
					Position: "Head Security",
				},
				Title:          "Pengajuan Penambahan Karyawan",
				SubmissionType: "Man Power",
				Attachment:     "facebook.com",
			},
		}

		repo.On("GetAllCc", succesID).Return(repoReturn, nil).Once()

		result, err := ccl.GetAllCcLogic(succesID)

		assert.NotEmpty(t, result)
		assert.Equal(t, repoReturn, ccResult)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}
