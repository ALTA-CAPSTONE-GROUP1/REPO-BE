package helper

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(fileContents *multipart.File, path string) ([]string, error) {
	var urls []string
	uploadResult, err := uploadFile(fileContents, path)
	if err != nil {
		return nil, err
	}
	urls = append(urls, uploadResult.SecureURL)
	return urls, nil
}

func uploadFile(content *multipart.File, path string) (*uploader.UploadResult, error) {
	cld, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryApiKey, config.CloudinaryApiScret)
	if err != nil {
		return nil, err
	}
	overwrite := true
	useFileName := true
	useFileNameDisplay := true
	uploadParams := uploader.UploadParams{
		PublicID:                 "epropProject",
		Folder:                   config.CloudinaryUploadFolder + path,
		UseFilename:              &useFileName,
		Overwrite:                &overwrite,
		UseFilenameAsDisplayName: &useFileNameDisplay,
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), *content, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("error in uploadin file %w", err)
	}
	return uploadResult, nil
}

// func AddStamp(signName string, currentLink string, path string)(string, error){

// }
