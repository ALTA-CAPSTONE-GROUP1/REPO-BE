package usecase

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/gommon/log"
)

type submissionLogic struct {
	sl submission.Repository
	u  helper.UploadInterface
}

func New(sr submission.Repository, u helper.UploadInterface) submission.UseCase {
	return &submissionLogic{
		sl: sr,
		u:  u,
	}
}

func (sr *submissionLogic) FindRequirementLogic(userID string, typeName string, value int) (submission.Core, error) {
	result, err := sr.sl.FindRequirement(userID, typeName, value)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return submission.Core{}, fmt.Errorf("data not found %w", err)
		} else if strings.Contains(err.Error(), "syntax") {
			return submission.Core{}, fmt.Errorf("internal server error %w", err)
		} else {
			return submission.Core{}, fmt.Errorf("unexpected error %w", err)
		}
	}

	return result, nil
}

func (sr *submissionLogic) AddSubmissionLogic(newSub submission.AddSubmissionCore, subFile *multipart.FileHeader) error {
	uploadUrl, err := sr.u.UploadFile(subFile, "/"+newSub.OwnerID)
	if err != nil {
		log.Errorf("error from third party upload file %w", err)
		return err
	}
	newSub.AttachmentLink = uploadUrl[0]
	newSub.Attachment = subFile.Filename

	if err := sr.sl.InsertSubmission(newSub); err != nil {
		log.Errorf("error on insert submission %w", err)
		if strings.Contains(err.Error(), "record not found") {
			return errors.New("record not found")
		}
		if strings.Contains(err.Error(), "duplicate") {
			log.Errorf("submission title or file is duplicate %w", err)
			return errors.New("record ")
		}
		if strings.Contains(err.Error(), "syntax") {
			return errors.New("syntax error")
		}
		return errors.New("unexpected error on inserting data")
	}

	return nil
}

func (sr *submissionLogic) UpdateDataByOwnerLogic(submission submission.UpdateCore, subFile *multipart.FileHeader) error {

	exist := sr.sl.FindFileData(submission.SubmissionID, subFile.Filename)

	if exist {
		log.Error("file name duplicate for update submissions")
		return errors.New("cannot upload same file and same file name to revise duplicate")
	}

	uploadUrl, err := sr.u.UploadFile(subFile, "/"+submission.UserID)
	if err != nil {
		log.Errorf("error fron third party upload file %w", err)
		return err
	}
	if len(uploadUrl) > 0 {
		submission.AttachmentLink = uploadUrl[0]
		submission.AttachmentName = subFile.Filename
	}
	err = sr.sl.UpdateDataByOwner(submission)
	if err != nil {
		if strings.Contains(err.Error(), "submission data not found") {
			return errors.New("submission data not found")
		}
		if strings.Contains(err.Error(), "status not") {
			return errors.New("submission status not sent")
		}
		if strings.Contains(err.Error(), "syntax") {
			return errors.New("internal server error(syntax)")
		}
		log.Errorf("unexpected error %w", err)
		return errors.New("server error, unexpected error")
	}

	return nil
}

func (sr *submissionLogic) GetAllSubmissionLogic(userID string, pr submission.GetAllQueryParams) ([]submission.AllSubmiisionCore, []submission.SubTypeChoices, error) {
	allsubmission, typelist, err := sr.sl.SelectAllSubmissions(userID, pr)
	if err != nil {
		log.Errorf("error on get all submission data", err)
		if strings.Contains(err.Error(), "record not found") {
			return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, errors.New("record not found")
		}
		if strings.Contains(err.Error(), "syntax") {
			return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, errors.New("syntax error")
		}
		return []submission.AllSubmiisionCore{}, []submission.SubTypeChoices{}, errors.New("unexpected error on inserting data")
	}

	return allsubmission, typelist, nil
}

func (sr *submissionLogic) GetSubmissionByIDLogic(submissionID int, userId string) (submission.GetSubmissionByIDCore, error) {
	result, err := sr.sl.SelectSubmissionByID(submissionID, userId)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return submission.GetSubmissionByIDCore{}, errors.New("record not found")
		}

		if strings.Contains(err.Error(), "syntax") {
			return submission.GetSubmissionByIDCore{}, errors.New("syntax error")
		}

		log.Errorf("unexpected error on getsubmissionByID")
		return submission.GetSubmissionByIDCore{}, fmt.Errorf("unexpected error %w", err)
	}

	return result, nil
}

func (sr *submissionLogic) DeleteSubmissionLogic(submissionID int, userID string) error {
	if err := sr.sl.DeleteSubmissionByID(submissionID, userID); err != nil {
		log.Errorf("error on calling delete submission ID")
		if strings.Contains(err.Error(), "not found") {
			return errors.New("data not found")
		}
		if strings.Contains(err.Error(), "sent") {
			return errors.New("unauthorized submission status is sent")
		}
		log.Errorf("unexpected error %w", err)
		return fmt.Errorf("unexppected error %w", err)
	}
	return nil
}
