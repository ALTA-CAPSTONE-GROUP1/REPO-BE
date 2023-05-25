package helper

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ALTA-CAPSTONE-GROUP1/e-proposal-BE/app/config"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/labstack/gommon/log"
	"github.com/phpdave11/gofpdf"
	fpdi "github.com/phpdave11/gofpdf/contrib/gofpdi"
)

func UpdateFile(currentLink string, approverName string, approverPosition string, subTitle string, signName string, path string) (string, string, error) {
	msgBody := fmt.Sprintf(`this message us autogenerate from epropApp this submission are approved by %s, %s,
	SignID = %s`, approverName, approverPosition, signName)
	outputpdf := "helper/output.pdf"
	createdPdf := CreatePDF(subTitle, msgBody, outputpdf)

	downloadedPdf := "helper/downloaded.pdf"

	err := downloadFile(currentLink, downloadedPdf)
	if err != nil {
		log.Errorf("error on downoading cloudinary file %w", err)
		return "", "", err
	}
	mergedFiles := "mergedfiles.pdf"
	err = mergePDFs(mergedFiles, downloadedPdf, createdPdf)
	if err != nil {
		log.Errorf("error on merging pdf %w", err)
		return "", "", err
	}
	fmt.Println("merged file berhasil dibuat")

	newUrl, err := UploadNewData("./mergedfiles.pdf", "/"+approverPosition)
	if err != nil {
		log.Errorf("error on upload pdf %s", err.Error())
		return "", "", err
	}
	fmt.Println("SAMPAI AFTER UPLOAD")
	err = os.Remove("./helper/downloaded.pdf")
	if err != nil {
		log.Errorf("error on on removing file created %w", err)
	}
	fmt.Println("SAMPAI AFTER REMOVE")
	err = os.Remove("./helper/output.pdf")
	if err != nil {
		log.Errorf("error on remove downloadedPdf %w", err)
	}
	err = os.Remove("./mergedfiles.pdf")
	if err != nil {
		log.Errorf("error on remove downloadedPdf %w", err)
	}
	fmt.Println("SAMPAI AFTER REMOVE merge")
	newName, err := GenerateUniqueSign(signName)
	if err != nil{
		log.Errorf(err.Error())
	}
	return newName, newUrl, nil
}

func CreatePDF(subTitle string, msgBody string, path string) string {
	pdf := gofpdf.New("P", "mm", "A4", "")

	pdf.AddPage()

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, subTitle)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 10, msgBody, "", "", false)
	pdf.Ln(12)

	pdf.SetFont("Arial", "I", 10)
	footerStr := fmt.Sprintf("Page %d From EpropApp | Date: %s", pdf.PageNo(), time.Now().Add(7*time.Hour).Format("02 January 2006"))
	pdf.SetY(-15)
	pdf.CellFormat(0, 10, footerStr, "", 0, "C", false, 0, "")

	pdf.AliasNbPages("")

	err := pdf.OutputFileAndClose(path)
	if err != nil {
		log.Errorf("error on creating pdf %w", err)
		return ""
	}

	fmt.Println("File PDF berhasil dibuat!")
	pdf.Close()
	return path
}

func downloadFile(url, downloadPath string) error {
	response, err := http.Get(url)
	if err != nil {
		log.Errorf("error on getting cloudinary file")
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(downloadPath)
	if err != nil {
		log.Errorf("error on creating downloadedPath %w", err)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Errorf("error on copying file downloaded to server %w", err)
		return err
	}

	return nil
}

func mergePDFs(destMerge string, files ...string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")

	for _, file := range files {
		fmt.Println("++++++++++Perulangan merge+++++++++++++")
		importedFile := fpdi.ImportPage(pdf, file, 1, "/MediaBox")
		pdf.AddPage()
		pdf.SetFont("Arial", "", 12)
		fpdi.UseImportedTemplate(pdf, importedFile, 20, 50, 150, 0)
	}
	fmt.Println("++++++++++++ Beres Perulangan +++++++`")
	err := pdf.OutputFileAndClose(destMerge)
	if err != nil {
		log.Errorf("error on creating merged file", err)
		return err
	}
	pdf.Close()
	return nil
}

func UploadNewData(filePath string, cloudinaryPath string) (string, error) {
	cld, err := cloudinary.NewFromParams(config.CloudinaryName, config.CloudinaryApiKey, config.CloudinaryApiScret)
	if err != nil {
		return "", err
	}

	overwrite := true
	useFileName := true
	useFileNameDisplay := true

	UploadParams := uploader.UploadParams{
		PublicID:                 "epropProject",
		Folder:                   config.CloudinaryUploadFolder + cloudinaryPath,
		UseFilename:              &useFileName,
		Overwrite:                &overwrite,
		UseFilenameAsDisplayName: &useFileNameDisplay,
	}

	resp, err := cld.Upload.Upload(context.Background(), filePath, UploadParams)
	if err != nil {
		log.Errorf("error on uploading new pdf %s", err.Error())
		return "", err
	}

	return resp.SecureURL, err
}
