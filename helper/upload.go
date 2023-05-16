package helper

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadFile(fileContents interface{}, path string) ([]string, error) {
	var urls []string
	switch cnv := fileContents.(type) {
	case []*multipart.File:
		for _, content := range cnv {
			uploadResult, err := uploadFile(content, path)
			if err != nil {
				return nil, err
			}
			urls = append(urls, uploadResult.SecureURL)
		}
	case *multipart.File:
		uploadResult, err := uploadFile(cnv, path)
		if err != nil {
			return nil, err
		}
		urls = append(urls, uploadResult.SecureURL)
		fmt.Println(urls)
		return urls, nil
	}
	fmt.Println(urls)
	return urls, nil
}

func uploadFile(content *multipart.File, path string) (*uploader.UploadResult, error) {
	cld, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryApiKey, config.CloudinaryApiScret)
	if err != nil {
		return nil, err
	}

	uploadParams := uploader.UploadParams{
		Folder: config.CloudinaryUploadFolder + path,
	}
	if err != nil {
		return nil, err
	}

	uploadResult, err := cld.Upload.Upload(context.Background(), *content, uploadParams)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return uploadResult, nil
}
