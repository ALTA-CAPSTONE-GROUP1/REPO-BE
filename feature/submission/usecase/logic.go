package usecase

import (
	"errors"
	"fmt"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/admin"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/feature/submission"
	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/helper"
	"github.com/labstack/gommon/log"
)

type submissionLogic struct {
	sl submission.Repository
}

func New(sr submission.Repository) submission.UseCase {
	return &submissionLogic{
		sl: sr,
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
	if subFile != nil {
		subFile, err := subFile.Open()
		if err != nil {
			log.Errorf("error in open subfile file")
			uploadUrl, err := helper.UploadFile(&subFile, "/submissionfiles")
			if err != nil {
				log.Errorf("error fron third party upload file %w", err)
				return err
			}
			newSub.AttachmentLink = uploadUrl[0]
			newSub.Attachment = newSub.OwnerID + strconv.FormatInt(time.Now().Unix(), 10)
		}
	}

	if err := sr.sl.InsertSubmission(newSub); err != nil {
		log.Errorf("error on insert submission %w", err)
		if strings.Contains(err.Error(), "record not found") {
			return errors.New("record not found")
		}
		if strings.Contains(err.Error(), "syntax") {
			return errors.New("syntax error")
		}
		return errors.New("unexpected error on inserting data")
	}

	return nil
}

func (sr *submissionLogic) GetAllSubmissionLogic(userID string, pr submission.GetAllQueryParams) ([]submission.AllSubmiisionCore, []admin.Type, error) {
	allsubmission, typelist, err := sr.sl.SelectAllSubmissions(userID, pr)
	if err != nil {
		log.Errorf("error on get all submission data", err)
		if strings.Contains(err.Error(), "record not found") {
			return []submission.AllSubmiisionCore{}, []admin.Type{}, errors.New("record not found")
		}
		if strings.Contains(err.Error(), "syntax") {
			return []submission.AllSubmiisionCore{}, []admin.Type{}, errors.New("syntax error")
		}
		return []submission.AllSubmiisionCore{}, []admin.Type{}, errors.New("unexpected error on inserting data")
	}

	return allsubmission, typelist, nil
}
