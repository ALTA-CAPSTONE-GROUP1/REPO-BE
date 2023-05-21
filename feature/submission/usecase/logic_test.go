package usecase_test

import (
	"mime/multipart"

	"os"
	"testing"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/mocks"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission/usecase"
	helperMocks "github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper/mocks"
	"github.com/stretchr/testify/assert"
)

// func TestFindRequirementLogic(t *testing.T) {
// 	repo := mocks.NewRepository(t)
// 	ul := usecase.New(repo)

// 	t.Run("Success Get Requirements", func(t *testing.T) {

// 		repo.On("FindRequirement", "RSM1", "Pengadaan", 30000000).
// 			Return(submission.Core{
// 				To: []submission.ToApprover{
// 					{
// 						ApproverPosition: "National Sales Manager",
// 						ApproverId:       "NSM2",
// 						ApproverName:     "Bohang",
// 					}, {
// 						ApproverPosition: "National Marketing Manager",
// 						ApproverId:       "NMM5",
// 						ApproverName:     "Didiek",
// 					},
// 				},
// 				CC: []submission.CcApprover{
// 					{
// 						CcPosition: "Regional Administration Manager",
// 						CcName:     "Puji Astuti",
// 						CcId:       "RAM78",
// 					}, {
// 						CcPosition: "Nasional Adminstration Manager",
// 						CcName:     "Puja Astuta",
// 						CcId:       "NAM12",
// 					},
// 				},
// 				Requirement: "Foto KTP, Penawaran pembanding",
// 			}, nil).Once()

// 		result, err := ul.FindRequirementLogic("RSM1", "Pengadaan", 30000000)

// 		assert.NoError(t, err)
// 		assert.NotNil(t, result)
// 		assert.Equal(t, result, submission.Core{
// 			To: []submission.ToApprover{
// 				{
// 					ApproverPosition: "National Sales Manager",
// 					ApproverId:       "NSM2",
// 					ApproverName:     "Bohang",
// 				}, {
// 					ApproverPosition: "National Marketing Manager",
// 					ApproverId:       "NMM5",
// 					ApproverName:     "Didiek",
// 				},
// 			},
// 			CC: []submission.CcApprover{
// 				{
// 					CcPosition: "Regional Administration Manager",
// 					CcName:     "Puji Astuti",
// 					CcId:       "RAM78",
// 				}, {
// 					CcPosition: "Nasional Adminstration Manager",
// 					CcName:     "Puja Astuta",
// 					CcId:       "NAM12",
// 				},
// 			},
// 			Requirement: "Foto KTP, Penawaran pembanding",
// 		})

// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Data Not Found", func(t *testing.T) {
// 		repo.On("FindRequirement", "ASM2", "Pengadaan", 60000000).
// 			Return(submission.Core{}, errors.New("record not found")).Once()

// 		result, err := ul.FindRequirementLogic("ASM2", "Pengadaan", 60000000)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "data not found record not found")
// 		assert.Equal(t, submission.Core{}, result)

// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Internal Server Error", func(t *testing.T) {
// 		repo.On("FindRequirement", "AMM3", "Penambahan Karyawan", 50000).
// 			Return(submission.Core{}, errors.New("syntax error")).Once()

// 		result, err := ul.FindRequirementLogic("AMM3", "Penambahan Karyawan", 50000)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "internal server error syntax error")
// 		assert.Equal(t, submission.Core{}, result)

// 		repo.AssertExpectations(t)
// 	})

// 	t.Run("Unexpected Error", func(t *testing.T) {
// 		repo.On("FindRequirement", "ASSpv1", "Trade Promo", 2500000).
// 			Return(submission.Core{}, errors.New("unexpected error")).Once()

// 		result, err := ul.FindRequirementLogic("ASSpv1", "Trade Promo", 2500000)

// 		assert.Error(t, err)
// 		assert.EqualError(t, err, "unexpected error unexpected error")
// 		assert.Equal(t, submission.Core{}, result)

// 		repo.AssertExpectations(t)
// 	})
// }

//proses pembukaan datanya di helpernya saja

func TestAddSubmissionLogic(t *testing.T) {
	repo := mocks.NewRepository(t)
	upload := helperMocks.NewUploadInterface(t)
	ul := usecase.New(repo, upload)

	t.Run("Succes Add Submission", func(t *testing.T) {
		file, err := os.Open("./example2.txt")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		subFile := &multipart.FileHeader{
			Filename: "example2.txt",
		}

		upload.On("UploadFile", subFile, "/NSM1").Return([]string{"example.com/sadada"}, nil).Once()
		repo.On("InsertSubmission", submission.AddSubmissionCore{
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
		}).Return(nil).Once()

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

		err = ul.AddSubmissionLogic(newSub, subFile)

		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	// t.Run("Record Not Found", func(t *testing.T) {
	// 	repo.On("InsertSubmission", submission.AddSubmissionCore{
	// 		OwnerID: "NSM1",
	// 		ToApprover: []submission.ToApprover{
	// 			{
	// 				ApproverPosition: "Reg Man",
	// 				ApproverId:       "RM1",
	// 				ApproverName:     "Kholil",
	// 			},
	// 		},
	// 		CC: []submission.CcApprover{
	// 			{
	// 				CcPosition: "AdMin Nasional",
	// 				CcName:     "Bani",
	// 				CcId:       "AN1",
	// 			},
	// 		},
	// 		SubmissionType:    "Memo Credit",
	// 		SubmissiontTypeID: 1,
	// 		Status:            "Sent",
	// 		SubmissionValue:   50000000,
	// 		Title:             "Berita Acara Credit Vendor A",
	// 		Message:           "Berikut saya lampirkan berita acara credit",
	// 		Attachment:        "file.txt",
	// 		AttachmentLink:    "res.cloudinary.com/aiueo",
	// 	}).Return(errors.New("record not found")).Once()

	// 	newSub := submission.AddSubmissionCore{
	// 		OwnerID: "NSM1",
	// 		ToApprover: []submission.ToApprover{
	// 			{
	// 				ApproverPosition: "Reg Man",
	// 				ApproverId:       "RM1",
	// 				ApproverName:     "Kholil",
	// 			},
	// 		},
	// 		CC: []submission.CcApprover{
	// 			{
	// 				CcPosition: "AdMin Nasional",
	// 				CcName:     "Bani",
	// 				CcId:       "AN1",
	// 			},
	// 		},
	// 		SubmissionType:    "Memo Credit",
	// 		SubmissiontTypeID: 1,
	// 		Status:            "Sent",
	// 		SubmissionValue:   50000000,
	// 		Title:             "Berita Acara Credit Vendor A",
	// 		Message:           "Berikut saya lampirkan berita acara credit",
	// 		Attachment:        "file.txt",
	// 		AttachmentLink:    "res.cloudinary.com/aiueo",
	// 	}

	// 	err := ul.AddSubmissionLogic(newSub, subFile)

	// 	assert.ErrorContains(t, err, "not found")
	// 	repo.AssertExpectations(t)
	// })

	// t.Run("Record Not Found", func(t *testing.T) {
	// 	repo.On("InsertSubmission", submission.AddSubmissionCore{
	// 		OwnerID: "NSM1",
	// 		ToApprover: []submission.ToApprover{
	// 			{
	// 				ApproverPosition: "Reg Man",
	// 				ApproverId:       "RM1",
	// 				ApproverName:     "Kholil",
	// 			},
	// 		},
	// 		CC: []submission.CcApprover{
	// 			{
	// 				CcPosition: "AdMin Nasional",
	// 				CcName:     "Dani",
	// 				CcId:       "JN1",
	// 			},
	// 		},
	// 		SubmissionType:    "Mamo Credit",
	// 		SubmissiontTypeID: 1,
	// 		Status:            "Sent",
	// 		SubmissionValue:   50000000,
	// 		Title:             "Berita Acara Credit Vendor A",
	// 		Message:           "Berikut saya lampirkan berita acara credit",
	// 		Attachment:        "file.txt",
	// 		AttachmentLink:    "res.cloudinary.com/aiueo",
	// 	}).Return(errors.New("syntax")).Once()

	// 	newSub := submission.AddSubmissionCore{
	// 		OwnerID: "NSM1",
	// 		ToApprover: []submission.ToApprover{
	// 			{
	// 				ApproverPosition: "Reg Man",
	// 				ApproverId:       "RM1",
	// 				ApproverName:     "Kholil",
	// 			},
	// 		},
	// 		CC: []submission.CcApprover{
	// 			{
	// 				CcPosition: "AdMin Nasional",
	// 				CcName:     "Dani",
	// 				CcId:       "JN1",
	// 			},
	// 		},
	// 		SubmissionType:    "Mamo Credit",
	// 		SubmissiontTypeID: 1,
	// 		Status:            "Sent",
	// 		SubmissionValue:   50000000,
	// 		Title:             "Berita Acara Credit Vendor A",
	// 		Message:           "Berikut saya lampirkan berita acara credit",
	// 		Attachment:        "file.txt",
	// 		AttachmentLink:    "res.cloudinary.com/aiueo",
	// 	}
	// 	subFile := &multipart.FileHeader{
	// 		Filename: "example3",
	// 		Header:   make(textproto.MIMEHeader),
	// 		Size:     repo.TestData().Value().Int64(),
	// 	}
	// 	err := ul.AddSubmissionLogic(newSub, subFile)

	// 	assert.ErrorContains(t, err, "syntax")
	// 	repo.AssertExpectations(t)
	// })

	// t.Run("Record Not Found", func(t *testing.T) {
	// 	repo.On("InsertSubmission", submission.AddSubmissionCore{
	// 		OwnerID: "NSM1",
	// 		ToApprover: []submission.ToApprover{
	// 			{
	// 				ApproverPosition: "Area Sales Supervsior",
	// 				ApproverId:       "Asspv",
	// 				ApproverName:     "Agung",
	// 			},
	// 		},
	// 		CC: []submission.CcApprover{
	// 			{
	// 				CcPosition: "Finance Nasional",
	// 				CcName:     "joni",
	// 				CcId:       "FN",
	// 			},
	// 		},
	// 		SubmissionType:    "Memo Credit",
	// 		SubmissiontTypeID: 1,
	// 		Status:            "Sent",
	// 		SubmissionValue:   50000000,
	// 		Title:             "Berita Acara Credit Vendor A",
	// 		Message:           "Berikut saya lampirkan berita acara credit",
	// 		Attachment:        "file.txt",
	// 		AttachmentLink:    "res.cloudinary.com/aiueo",
	// 	}).Return(errors.New("another error")).Once()

	// 	newSub := submission.AddSubmissionCore{
	// 		OwnerID: "NSM1",
	// 		ToApprover: []submission.ToApprover{
	// 			{
	// 				ApproverPosition: "Area Sales Supervsior",
	// 				ApproverId:       "Asspv",
	// 				ApproverName:     "Agung",
	// 			},
	// 		},
	// 		CC: []submission.CcApprover{
	// 			{
	// 				CcPosition: "Finance Nasional",
	// 				CcName:     "joni",
	// 				CcId:       "FN",
	// 			},
	// 		},
	// 		SubmissionType:    "Memo Credit",
	// 		SubmissiontTypeID: 1,
	// 		Status:            "Sent",
	// 		SubmissionValue:   50000000,
	// 		Title:             "Berita Acara Credit Vendor A",
	// 		Message:           "Berikut saya lampirkan berita acara credit",
	// 		Attachment:        "file.txt",
	// 		AttachmentLink:    "res.cloudinary.com/aiueo",
	// 	}
	// 	subFile := &multipart.FileHeader{
	// 		Filename: "example4",
	// 		Header:   make(textproto.MIMEHeader),
	// 		Size:     repo.TestData().Value().Int64(),
	// 	}
	// 	err := ul.AddSubmissionLogic(newSub, subFile)

	// 	assert.ErrorContains(t, err, "unexpected")
	// 	repo.AssertExpectations(t)
	// })

}
