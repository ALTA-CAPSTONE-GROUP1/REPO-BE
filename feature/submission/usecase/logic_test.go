package usecase_test

import (
	"errors"
	"mime/multipart"

	"os"
	"testing"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/usecase"
	helperMocks "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper/mocks"
	"github.com/stretchr/testify/assert"
)

//proses pembukaan datanya di helpernya saja

func TestAddSubmissionLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	upload := helperMocks.NewUploadInterface(t)
	ul := usecase.New(repo, upload)

	t.Run("Succes Add Submission", func(t *testing.T) {
		fileSucces, _ := os.Open("./example2.txt")
		defer fileSucces.Close()
		subFile := &multipart.FileHeader{
			Filename: "example2.txt",
		}

		upload.On("UploadFile", subFile, "/NSM1").Return([]string{"example.com/sadada"}, nil).Once()
		succesCore := submission.AddSubmissionCore{
			OwnerID: "NSM1",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM1",
					ApproverName:     "Kholil",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "AdMin Nasional",
					CcName:     "Vona",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   8000000,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
			Attachment:        "example2.txt",
			AttachmentLink:    "example.com/sadada",
		}
		repo.On("InsertSubmission", succesCore).Return(nil).Once()

		newSub := submission.AddSubmissionCore{
			OwnerID: "NSM1",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM1",
					ApproverName:     "Kholil",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "AdMin Nasional",
					CcName:     "Vona",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   8000000,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
		}

		err := ul.AddSubmissionLogic(newSub, subFile)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Error from Third Party", func(t *testing.T) {

		failFile, _ := os.Open("./example3.txt")
		defer failFile.Close()
		failMulti := &multipart.FileHeader{
			Filename: "example3.txt",
		}
		failSubThirParty := submission.AddSubmissionCore{
			OwnerID: "AMM1",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM1",
					ApproverName:     "Kholal",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "AdMin Nasional",
					CcName:     "Vona",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   12231314,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
		}
		upload.On("UploadFile", failMulti, "/AMM1").Return([]string{""}, errors.New("anything")).Once()

		err := ul.AddSubmissionLogic(failSubThirParty, failMulti)

		assert.Error(t, err)
		assert.EqualError(t, err, "anything")
		repo.AssertExpectations(t)
	})

	t.Run("Error Record Not Found", func(t *testing.T) {

		failFile, _ := os.Open("./example4.txt")
		defer failFile.Close()
		multiRecordNotFound := &multipart.FileHeader{
			Filename: "example4.txt",
		}
		FailSubRecordLogic := submission.AddSubmissionCore{
			OwnerID: "Asspv1",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM1",
					ApproverName:     "Kholal",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "AdMin Nasional",
					CcName:     "Vona",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   12231314,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
			Attachment:        "example4.txt",
			AttachmentLink:    "example4.com/sadada",
		}
		failRecord := submission.AddSubmissionCore{
			OwnerID: "Asspv1",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM1",
					ApproverName:     "Kholal",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "AdMin Nasional",
					CcName:     "Vona",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   12231314,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
			Attachment:        "example4.txt",
			AttachmentLink:    "example4.com/sadada",
		}
		upload.On("UploadFile", multiRecordNotFound, "/Asspv1").Return([]string{"example4.com/sadada"}, nil).Once()
		repo.On("InsertSubmission", failRecord).Return(errors.New("record not found")).Once()

		err := ul.AddSubmissionLogic(FailSubRecordLogic, multiRecordNotFound)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "record not found")
		repo.AssertExpectations(t)
	})

	t.Run("Error Syntax", func(t *testing.T) {

		failFile, _ := os.Open("./example5.txt")
		defer failFile.Close()
		multiSyntaxFile := &multipart.FileHeader{
			Filename: "example5.txt",
		}
		failSyntaxLogic := submission.AddSubmissionCore{
			OwnerID: "WSS5",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM2",
					ApproverName:     "Hana",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "Jabatan Apa Saja",
					CcName:     "Vani",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   12231314,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
		}
		failSyntaxRepo := submission.AddSubmissionCore{
			OwnerID: "WSS5",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM2",
					ApproverName:     "Hana",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "Jabatan Apa Saja",
					CcName:     "Vani",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   12231314,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
			Attachment:        "example5.txt",
			AttachmentLink:    "example5.com/sadada",
		}
		upload.On("UploadFile", multiSyntaxFile, "/WSS5").Return([]string{"example5.com/sadada"}, nil).Once()
		repo.On("InsertSubmission", failSyntaxRepo).Return(errors.New("syntax")).Once()

		err := ul.AddSubmissionLogic(failSyntaxLogic, multiSyntaxFile)

		assert.Error(t, err)
		assert.EqualError(t, err, "syntax error")
		repo.AssertExpectations(t)
	})

	t.Run("Unexpec error", func(t *testing.T) {

		failFile, _ := os.Open("./example6.txt")
		defer failFile.Close()
		multiSyntaxFile := &multipart.FileHeader{
			Filename: "example6.txt",
		}
		FailUnexpectedLogic := submission.AddSubmissionCore{
			OwnerID: "AMTO",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM2",
					ApproverName:     "Hana",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "Jabatan Apa Saja",
					CcName:     "Vani",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   12231314,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
		}
		FailUnexpectedRepo := submission.AddSubmissionCore{
			OwnerID: "AMTO",
			ToApprover: []submission.ToApprover{
				{
					ApproverPosition: "Reg Man",
					ApproverId:       "RM2",
					ApproverName:     "Hana",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "Jabatan Apa Saja",
					CcName:     "Vani",
					CcId:       "AN1",
				},
			},
			SubmissionType:    "Memo Credit",
			SubmissiontTypeID: 1,
			Status:            "Sent",
			SubmissionValue:   12231314,
			Title:             "Berita Acara Credit Vendor A",
			Message:           "Berikut saya lampirkan berita acara credit",
			Attachment:        "example6.txt",
			AttachmentLink:    "example6.com/sadada",
		}
		upload.On("UploadFile", multiSyntaxFile, "/AMTO").Return([]string{"example6.com/sadada"}, nil).Once()
		repo.On("InsertSubmission", FailUnexpectedRepo).Return(errors.New("anything")).Once()

		err := ul.AddSubmissionLogic(FailUnexpectedLogic, multiSyntaxFile)

		assert.Error(t, err)
		assert.EqualError(t, err, "unexpected error on inserting data")
		repo.AssertExpectations(t)
	})
}

func TestFindRequirementLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	upload := helperMocks.NewUploadInterface(t)
	ul := usecase.New(repo, upload)

	t.Run("Record Not Found", func(t *testing.T) {
		repo.On("FindRequirement", "NAM1", "Pengadaan", 222000).
			Return(submission.Core{}, errors.New("record not found")).Once()
		resultNotFound, errorNotFound := ul.FindRequirementLogic("NAM1", "Pengadaan", 222000)

		assert.Error(t, errorNotFound)
		assert.EqualError(t, errorNotFound, "data not found record not found")
		assert.Empty(t, resultNotFound)
		repo.AssertExpectations(t)

	})

	t.Run("Internal Server Error", func(t *testing.T) {
		repo.On("FindRequirement", "NEM2", "Pengadaan", 202100).
			Return(submission.Core{}, errors.New("syntax error")).Once()
		result, err := ul.FindRequirementLogic("NEM2", "Pengadaan", 202100)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal server error syntax error")
		assert.Empty(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		repo.On("FindRequirement", "the2", "penawaran", 50200).
			Return(submission.Core{}, errors.New("unexpected error"))
		result, err := ul.FindRequirementLogic("the2", "penawaran", 50200)

		assert.Error(t, err)
		assert.EqualError(t, err, "unexpected error unexpected error")
		assert.Empty(t, result)
		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expectedResult := submission.Core{
			To: []submission.ToApprover{
				{
					ApproverPosition: "Nasional Manager",
					ApproverId:       "NAM@",
					ApproverName:     "Siapa sAja",
				},
			},
			CC: []submission.CcApprover{
				{
					CcPosition: "Nasional Admin Manager",
					CcName:     "Bulalak",
					CcId:       "Dia",
				},
			},
			Requirement: "KTP",
		}
		repo.On("FindRequirement", "NAM1", "Pengadaan", 2000).
			Return(expectedResult, nil)
		result, err := ul.FindRequirementLogic("NAM1", "Pengadaan", 2000)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		repo.AssertExpectations(t)
	})

}

func TestGetAllSubmissionLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	upload := helperMocks.NewUploadInterface(t)
	ul := usecase.New(repo, upload)

	t.Run("Record Not Found", func(t *testing.T) {
		queryParams := submission.GetAllQueryParams{}
		repo.On("SelectAllSubmissions", "NAM2", queryParams).
			Return([]submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, errors.New("record not found"))
		allSubmissions, typeList, err := ul.GetAllSubmissionLogic("NAM2", queryParams)

		assert.Error(t, err)
		assert.EqualError(t, err, "record not found")
		assert.Empty(t, allSubmissions)
		assert.Empty(t, typeList)
		repo.AssertExpectations(t)
	})

	t.Run("Syntax Error", func(t *testing.T) {
		queryParams := submission.GetAllQueryParams{}
		repo.On("SelectAllSubmissions", "ASSPV1", queryParams).
			Return([]submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, errors.New("syntax error"))
		allSubmissions, typeList, err := ul.GetAllSubmissionLogic("ASSPV1", queryParams)

		assert.Error(t, err)
		assert.EqualError(t, err, "syntax error")
		assert.Empty(t, allSubmissions)
		assert.Empty(t, typeList)
		repo.AssertExpectations(t)
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		queryParams := submission.GetAllQueryParams{}
		repo.On("SelectAllSubmissions", "NANA1", queryParams).
			Return([]submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, errors.New("unexpected error"))
		allSubmissions, typeList, err := ul.GetAllSubmissionLogic("NANA1", queryParams)

		assert.Error(t, err)
		assert.EqualError(t, err, "unexpected error on inserting data")
		assert.Empty(t, allSubmissions)
		assert.Empty(t, typeList)
		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expectedAllSubmissions := []submission.AllSubmiisionCore{}
		expectedTypeList := []submission.SubTypeChoices{}
		queryParams := submission.GetAllQueryParams{}
		repo.On("SelectAllSubmissions", "STPM", queryParams).
			Return(expectedAllSubmissions, expectedTypeList, nil)
		allSubmissions, typeList, err := ul.GetAllSubmissionLogic("STPM", queryParams)

		assert.NoError(t, err)
		assert.Equal(t, expectedAllSubmissions, allSubmissions)
		assert.Equal(t, expectedTypeList, typeList)
		repo.AssertExpectations(t)
	})
}

func TestGetSubmissionByIDLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	upload := helperMocks.NewUploadInterface(t)
	ul := usecase.New(repo, upload)

	t.Run("Record Not Found", func(t *testing.T) {
		submissionID := 123
		userID := "BAM"
		repo.On("SelectSubmissionByID", submissionID, userID).
			Return(submission.GetSubmissionByIDCore{}, errors.New("record not found"))
		result, err := ul.GetSubmissionByIDLogic(submissionID, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "record not found")
		assert.Equal(t, submission.GetSubmissionByIDCore{}, result)
		repo.AssertExpectations(t)
	})

	t.Run("Syntax Error", func(t *testing.T) {
		submissionID := 1
		userID := "NEM"
		repo.On("SelectSubmissionByID", submissionID, userID).
			Return(submission.GetSubmissionByIDCore{}, errors.New("syntax error"))
		result, err := ul.GetSubmissionByIDLogic(submissionID, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "syntax error")
		assert.Equal(t, submission.GetSubmissionByIDCore{}, result)
		repo.AssertExpectations(t)
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		submissionID := 22
		userID := "NIM"
		repo.On("SelectSubmissionByID", submissionID, userID).
			Return(submission.GetSubmissionByIDCore{}, errors.New("unexpected error"))
		result, err := ul.GetSubmissionByIDLogic(submissionID, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "unexpected error unexpected error")
		assert.Equal(t, submission.GetSubmissionByIDCore{}, result)
		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expectedResult := submission.GetSubmissionByIDCore{}
		submissionID := 1
		userID := "NOM"
		repo.On("SelectSubmissionByID", submissionID, userID).
			Return(expectedResult, nil)
		result, err := ul.GetSubmissionByIDLogic(submissionID, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedResult, result)
		repo.AssertExpectations(t)
	})
}

func TestDeleteSubmissionLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	upload := helperMocks.NewUploadInterface(t)
	ul := usecase.New(repo, upload)

	t.Run("Data Not Found", func(t *testing.T) {
		submissionID := 1
		userID := "NAM1"
		repo.On("DeleteSubmissionByID", submissionID, userID).
			Return(errors.New("not found"))
		err := ul.DeleteSubmissionLogic(submissionID, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "data not found")
		repo.AssertExpectations(t)
	})

	t.Run("Unauthorized Submission Status", func(t *testing.T) {
		submissionID := 2
		userID := "NAM2"
		repo.On("DeleteSubmissionByID", submissionID, userID).
			Return(errors.New("sent"))
		err := ul.DeleteSubmissionLogic(submissionID, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "unauthorized submission status is sent")
		repo.AssertExpectations(t)
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		submissionID := 3
		userID := "NAM3"
		repo.On("DeleteSubmissionByID", submissionID, userID).
			Return(errors.New("unexpected error"))
		err := ul.DeleteSubmissionLogic(submissionID, userID)

		assert.Error(t, err)
		assert.EqualError(t, err, "unexppected error unexpected error")
		repo.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		submissionID := 4
		userID := "NAM4"
		repo.On("DeleteSubmissionByID", submissionID, userID).
			Return(nil)
		err := ul.DeleteSubmissionLogic(submissionID, userID)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}

func TestUpdateDataByOwnerLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	upload := helperMocks.NewUploadInterface(t)
	ul := usecase.New(repo, upload)

	t.Run("Success Update Data by Owner", func(t *testing.T) {
		existingFile := &multipart.FileHeader{
			Filename: "existing21.txt",
		}
		submissionData := submission.UpdateCore{
			SubmissionID:   5556,
			UserID:         "user123",
			Message:        "Updated Message",
			AttachmentLink: "example.com/existing21",
			AttachmentName: "existing2.txt",
		}

		upload.On("UploadFile", existingFile, "/user123").Return([]string{}, nil).Once()
		repo.On("FindFileData", submissionData.SubmissionID, existingFile.Filename).Return(false).Once()
		repo.On("UpdateDataByOwner", submissionData).Return(nil).Once()

		err := ul.UpdateDataByOwnerLogic(submissionData, existingFile)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
		upload.AssertExpectations(t)
	})

	t.Run("Unexpected Error", func(t *testing.T) {
		existingFile := &multipart.FileHeader{
			Filename: "existing31.txt",
		}
		submissionData := submission.UpdateCore{
			SubmissionID:   1114,
			UserID:         "user456",
			Message:        "Updated Message",
			AttachmentLink: "example.com/existing",
			AttachmentName: "existing2.txt",
		}

		upload.On("UploadFile", existingFile, "/user456").Return([]string{}, nil).Once()
		repo.On("FindFileData", submissionData.SubmissionID, existingFile.Filename).Return(false).Once()
		repo.On("UpdateDataByOwner", submissionData).Return(errors.New("update failed")).Once()

		err := ul.UpdateDataByOwnerLogic(submissionData, existingFile)

		assert.Error(t, err)
		assert.EqualError(t, err, "server error, unexpected error")
		repo.AssertExpectations(t)
		upload.AssertExpectations(t)
	})
	t.Run("Syntax Error", func(t *testing.T) {
		existingFile := &multipart.FileHeader{
			Filename: "existing22131.txt",
		}
		submissionData := submission.UpdateCore{
			SubmissionID:   27,
			UserID:         "user870",
			Message:        "Updated Message",
			AttachmentLink: "example.com/existing22131",
			AttachmentName: "existing22131.txt",
		}

		upload.On("UploadFile", existingFile, "/user870").Return([]string{"example.com/existing22131"}, nil).Once()
		repo.On("FindFileData", submissionData.SubmissionID, existingFile.Filename).Return(false).Once()
		repo.On("UpdateDataByOwner", submissionData).Return(errors.New("syntax")).Once()

		err := ul.UpdateDataByOwnerLogic(submissionData, existingFile)

		assert.Error(t, err)
		assert.EqualError(t, err, "internal server error(syntax)")
		repo.AssertExpectations(t)
		upload.AssertExpectations(t)
	})

	t.Run("Error from Status Not Sent", func(t *testing.T) {
		existingFile := &multipart.FileHeader{
			Filename: "existing777.txt",
		}
		submissionData := submission.UpdateCore{
			SubmissionID:   2221,
			UserID:         "user09090",
			Message:        "Updated Message",
			AttachmentLink: "example.com/existing777",
			AttachmentName: "existing777.txt",
		}

		upload.On("UploadFile", existingFile, "/user09090").Return([]string{"example.com/existing777"}, nil).Once()
		repo.On("FindFileData", submissionData.SubmissionID, existingFile.Filename).Return(false).Once()
		repo.On("UpdateDataByOwner", submissionData).Return(errors.New("status not")).Once()

		err := ul.UpdateDataByOwnerLogic(submissionData, existingFile)

		assert.Error(t, err)
		assert.EqualError(t, err, "submission status not sent")
		repo.AssertExpectations(t)
		upload.AssertExpectations(t)
	})

	t.Run("Error from Submission Data Not Found", func(t *testing.T) {
		existingFile := &multipart.FileHeader{
			Filename: "existing35.txt",
		}
		submissionData := submission.UpdateCore{
			SubmissionID:   90,
			UserID:         "user12309",
			Message:        "Updated Message",
			AttachmentLink: "example.com/existing35",
			AttachmentName: "existing35.txt",
		}

		upload.On("UploadFile", existingFile, "/user12309").Return([]string{"example.com/existing35"}, nil).Once()
		repo.On("FindFileData", submissionData.SubmissionID, existingFile.Filename).Return(false).Once()
		repo.On("UpdateDataByOwner", submissionData).Return(errors.New("submission data not found")).Once()

		err := ul.UpdateDataByOwnerLogic(submissionData, existingFile)

		assert.Error(t, err)
		assert.EqualError(t, err, "submission data not found")
		repo.AssertExpectations(t)
		upload.AssertExpectations(t)
	})

	t.Run("Third Party Error", func(t *testing.T) {
		existingFile := &multipart.FileHeader{
			Filename: "existing226.txt",
		}
		submissionData := submission.UpdateCore{
			SubmissionID:   102,
			UserID:         "NAM233",
			Message:        "Updated Message",
			AttachmentLink: "example.com/existing2226",
			AttachmentName: "existing2226.txt",
		}

		repo.On("FindFileData", submissionData.SubmissionID, existingFile.Filename).Return(false).Once()
		upload.On("UploadFile", existingFile, "/NAM233").Return([]string{}, errors.New("TLS HANDSHAKE")).Once()

		err := ul.UpdateDataByOwnerLogic(submissionData, existingFile)

		assert.Error(t, err)
		assert.EqualError(t, err, "TLS HANDSHAKE")
		repo.AssertExpectations(t)
		upload.AssertExpectations(t)
	})

	t.Run("Error Duplicate", func(t *testing.T) {
		missingFile := &multipart.FileHeader{
			Filename: "missing.txt",
		}
		submissionNoData := submission.UpdateCore{
			SubmissionID:   212,
			UserID:         "user789",
			Message:        "Updated Message",
			AttachmentLink: "example.com/missing",
			AttachmentName: "missing.txt",
		}

		repo.On("FindFileData", submissionNoData.SubmissionID, missingFile.Filename).Return(true).Once()

		err := ul.UpdateDataByOwnerLogic(submissionNoData, missingFile)

		assert.Error(t, err)
		assert.EqualError(t, err, "cannot upload same file and same file name to revise")
		repo.AssertExpectations(t)
		upload.AssertExpectations(t)
	})
}
