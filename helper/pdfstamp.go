package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/phpdave11/gofpdf"
	fpdi "github.com/phpdave11/gofpdf/contrib/gofpdi"
)

type UpdateFile struct {
	u UploadInterface
}

type UpdateInterface interface {
	UpdateFile(link string, appName string, appPos string, subTitle string, signName string, path string) (string, []string, error)
}

func NewUpdateInterface(u UploadInterface) UpdateInterface {
	return &UpdateFile{
		u: u,
	}
}

func (uf *UpdateFile) UpdateFile(currentLink string, approverName string, approverPosition string, subTitle string, signName string, path string) (string, []string, error) {
	msgBody := fmt.Sprintf(`this message us autogenerate from epropApp this submission are approved by %s, %s,
	SignID = %s`, approverName, approverPosition, signName)

	createdPdf := CreatePDF(subTitle, msgBody)

	downloadedPdf := "downloaded.pdf"

	err := downloadFile(currentLink, downloadedPdf)
	if err != nil {
		log.Errorf("error on downoading cloudinary file %w", err)
		return "", []string{}, err
	}
	log.Info("File downloaded successfully!")

	mergedFiles := "mergedfiles.pdf"
	err = mergePDFs(mergedFiles, createdPdf, downloadedPdf)
	if err != nil {
		log.Errorf("error on merging pdf %w", err)
		return "", []string{}, err
	}
	file, err := os.Open(mergedFiles)
	if err != nil {
		log.Errorf("error on opening mergedfile %w", err)
		return "", []string{}, err
	}
	defer file.Close()

	fileHead := &multipart.FileHeader{
		Filename: file.Name(),
	}

	url, err := uf.u.UploadFile(fileHead, "/cobadulu")
	if err != nil {
		log.Errorf("error on calling upload file %w", err)
		return "", []string{}, err
	}

	err = os.Remove(mergedFiles)
	if err != nil {
		log.Errorf("error on remove mergedfile")
	}

	return file.Name(), url, nil
}

func CreatePDF(subTitle string, msgBody string) string {
	outputPath := "output.pdf"

	if _, err := os.Stat(outputPath); err != nil {
		err := os.Remove(outputPath)
		if err != nil {
			log.Errorf("error on deleting outputfile.pdf, %w", err)
		}
	}

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

	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		log.Errorf("error on creating pdf %w", err)
		return ""
	}

	fmt.Println("File PDF berhasil dibuat!")
	pdf.Close()
	return outputPath
}

func downloadFile(url, downloadPath string) error {

	if _, err := os.Stat(downloadPath); err != nil {
		err := os.Remove(downloadPath)
		if err != nil {
			log.Errorf("error on deleting outputfile.pdf, %w", err)
		}
	}

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
		importedFile := fpdi.ImportPage(pdf, file, 1, "/MediaBox")
		pdf.AddPage()
		pdf.SetFont("Arial", "", 12)
		fpdi.UseImportedTemplate(pdf, importedFile, 20, 50, 150, 0)
	}
	err := pdf.OutputFileAndClose(destMerge)
	if err != nil {
		log.Errorf("error on creating merged file", err)
		return err
	}

	return nil
}
